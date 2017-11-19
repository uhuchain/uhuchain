/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"fmt"
	"time"

	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/admin"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/chclient"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/uhuchain/uhu-hlf-client/log"
)

//ClientConfig holds the configuration of the blockchain interface
type ClientConfig struct {
	Client      fab.FabricClient
	SDK         *deffab.FabricSDK
	ConfigFile  string
	sdkUserUser ca.User
	Initialized bool
}

/*
Initialize the client for the blockchain based on a configuration file
Parameters:
 - orgID: Organization ID the clients belong to.
			Must be included in the clients configuration file
 - sdkUser: Name of a user belonging to the organization
*/
func (setup *ClientConfig) Initialize(orgID string, sdkUser string) error {

	sdkOptions := deffab.Options{
		ConfigFile: setup.ConfigFile,
	}

	sdk, err := deffab.NewSDK(sdkOptions)
	if err != nil {
		return err
	}
	setup.SDK = sdk

	session, err := sdk.NewPreEnrolledUserSession(orgID, sdkUser)
	if err != nil {
		return err
	}

	sc, err := sdk.NewSystemClient(session)
	if err != nil {
		return err
	}
	setup.Client = sc

	setup.sdkUserUser = session.Identity()

	setup.Initialized = true

	return nil
}

/*
SetupChannel initializes and returns a channel object.
Note: This function does not send the actual proposal to an orderer.
You have to use the CreateChannel function for that
*/
func (setup *ClientConfig) setupChannel(channelID string, orgs []string) (fab.Channel, error) {
	if !setup.Initialized {
		return nil, fmt.Errorf("config not initialzed")
	}

	channel, err := setup.Client.NewChannel(channelID)
	if err != nil {
		return nil, err
	}

	ordererConfig, err := setup.Client.Config().RandomOrdererConfig()

	if err != nil {
		return nil, err
	}
	serverHostOverride := ""
	if str, ok := ordererConfig.GRPCOptions["ssl-target-name-override"].(string); ok {
		serverHostOverride = str
	}
	orderer, err := orderer.NewOrderer(ordererConfig.URL, ordererConfig.TLSCACerts.Path,
		serverHostOverride, setup.Client.Config())
	if err != nil {
		return nil, err
	}
	err = channel.AddOrderer(orderer)
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		peerConfig, err := setup.Client.Config().PeersConfig(org)
		if err != nil {
			return nil, err
		}
		for _, p := range peerConfig {
			serverHostOverride = ""
			if str, ok := p.GRPCOptions["ssl-target-name-override"].(string); ok {
				serverHostOverride = str
			}
			endorser, err := deffab.NewPeer(p.URL, p.TLSCACerts.Path,
				serverHostOverride, setup.Client.Config())
			if err != nil {
				return nil, err
			}
			err = channel.AddPeer(endorser)
			if err != nil {
				return nil, err
			}
		}
	}
	return channel, nil
}

/*
CreateChannel creates a channel based on a channel artifact.
All peers ascociated with a given organization in the clients config file
will be added as endorsers to channel object.
If the CreateChannelRequests was succesfull, it will automatically
join the client and return the channel object.
Remark: ChannelID needs to be the same as the channel name used creating
the artifact.
*/
func (setup *ClientConfig) CreateChannel(channelID string, ordererOrg string,
	orderersdkUserName string, orgs []string, channelArtifactPath string) (fab.Channel, error) {

	channel, err := setup.setupChannel(channelID, orgs)
	if err != nil {
		return nil, err
	}

	orderersdkUser, err := setup.SDK.NewPreEnrolledUser(ordererOrg, orderersdkUserName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed getting orderer sdkUser user")
	}

	// Check if primary peer has joined channel
	alreadyJoined, err := HasPrimaryPeerJoinedChannel(setup.Client, channel)
	if err != nil {
		return nil, errors.WithMessage(err, "failed while checking if primary peer has already joined channel")
	}

	if !alreadyJoined {
		// Create, initialize and join channel
		if err = admin.CreateOrUpdateChannel(setup.Client, orderersdkUser, setup.sdkUserUser,
			channel, channelArtifactPath); err != nil {
			return nil, errors.WithMessage(err, "CreateChannel failed")
		}
		time.Sleep(time.Second * 3)

		if err = channel.Initialize(nil); err != nil {
			return nil, errors.WithMessage(err, "channel init failed")
		}

		if err = admin.JoinChannel(setup.Client, setup.sdkUserUser, channel); err != nil {
			return nil, errors.WithMessage(err, "JoinChannel failed")
		}
	}
	return channel, nil
}

// GetEventHubForOrg initilizes the event hub for the peers of a give organization
func (setup *ClientConfig) GetEventHubForOrg(orgID string) (fab.EventHub, error) {
	eventHub, err := events.NewEventHub(setup.Client)
	if err != nil {
		return nil, err
	}
	foundEventHub := false
	peerConfigs, err := setup.Client.Config().PeersConfig(orgID)
	if err != nil {
		return nil, err
	}
	for _, p := range peerConfigs {
		if p.URL != "" {
			log.Info.Printf("EventHub connect to peer (%s)", p.URL)
			serverHostOverride := ""
			if str, ok := p.GRPCOptions["ssl-target-name-override"].(string); ok {
				serverHostOverride = str
			}
			eventHub.SetPeerAddr(p.EventURL, p.TLSCACerts.Path, serverHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, errors.WithMessage(err, "EventHub not found.")
	}

	return eventHub, nil
}

// RegisterTxEvent registers on the given eventhub for the give transaction
// returns a boolean channel which receives true when the event is complete
// and an error channel for errors
// TODO - Duplicate
func (setup *ClientConfig) RegisterTxEvent(txID apitxn.TransactionID,
	eventHub fab.EventHub) (chan bool, chan error) {
	done := make(chan bool)
	fail := make(chan error)

	eventHub.RegisterTxEvent(txID, func(txId string, errorCode pb.TxValidationCode, err error) {
		if err != nil {
			log.Error.Printf("Received error event for txid(%s)", txId)
			fail <- err
		} else {
			log.Error.Printf("Received success event for txid(%s)", txId)
			done <- true
		}
	})

	return done, fail
}

// InstallCC installs chaincode on a channel
func (setup *ClientConfig) InstallCC(channel fab.Channel, chainCodeID string,
	localChaincodeBasePath string, chainCodePath string, chainCodeVersion string,
	chaincodePackage []byte) error {

	if err := admin.SendInstallCC(setup.Client, chainCodeID, chainCodePath,
		chainCodeVersion, chaincodePackage, channel.Peers(),
		localChaincodeBasePath); err != nil {
		return errors.WithMessage(err, "SendInstallProposal failed")
	}
	log.Info.Printf("successfully send install proposal for chaincode %s", chainCodeID)
	return nil
}

// InstantiateCC ...
func (setup *ClientConfig) InstantiateCC(channel fab.Channel, chainCodeID string,
	chainCodePath string, chainCodeVersion string,
	args [][]byte, eventHub fab.EventHub) error {

	chaincodePolicy := cauthdsl.SignedByMspMember(setup.Client.UserContext().MspID())

	if err := admin.SendInstantiateCC(channel, chainCodeID, args,
		chainCodePath, chainCodeVersion, chaincodePolicy,
		[]apitxn.ProposalProcessor{channel.PrimaryPeer()}, eventHub); err != nil {
		return errors.WithMessage(err, "SendInstantiateProposal failed")
	}
	log.Info.Printf("successfully send instantiate proposal for chaincode %s", chainCodeID)
	return nil
}

// UpgradeCC ...
func (setup *ClientConfig) UpgradeCC(channel fab.Channel, chainCodeID string,
	chainCodePath string, chainCodeVersion string,
	args [][]byte, eventHub fab.EventHub) error {

	chaincodePolicy := cauthdsl.SignedByMspMember(setup.Client.UserContext().MspID())

	return admin.SendUpgradeCC(channel, chainCodeID, args, chainCodePath, chainCodeVersion, chaincodePolicy, []apitxn.ProposalProcessor{channel.PrimaryPeer()}, eventHub)
}

// SetupCC installs AND initiates a chaincode on a given channel
func (setup *ClientConfig) SetupCC(channel fab.Channel, chainCodeID string,
	localChaincodeBasePath string, chainCodePath string, chainCodeVersion string,
	chaincodePackage []byte, args [][]byte, eventHub fab.EventHub) error {
	err := setup.InstallCC(channel, chainCodeID, localChaincodeBasePath,
		chainCodePath, chainCodeVersion, chaincodePackage)
	if err != nil {
		return err
	}
	return setup.InstantiateCC(channel, chainCodeID, chainCodePath,
		chainCodeVersion, args, eventHub)
}

// SetupUpgradeCC installs AND upgrades a chaincode on a given channel
func (setup *ClientConfig) SetupUpgradeCC(channel fab.Channel, chainCodeID string,
	localChaincodeBasePath string, chainCodePath string, chainCodeVersion string,
	chaincodePackage []byte, args [][]byte, eventHub fab.EventHub) error {
	err := setup.InstallCC(channel, chainCodeID, localChaincodeBasePath,
		chainCodePath, chainCodeVersion, chaincodePackage)
	if err != nil {
		return err
	}
	return setup.UpgradeCC(channel, chainCodeID, chainCodePath,
		chainCodeVersion, args, eventHub)
}

// Query the state of an asset on a given channel
func (setup *ClientConfig) Query(channelID string, user string,
	chainCodeID string, args [][]byte) ([]byte, error) {
	chClient, err := setup.SDK.NewChannelClient(channelID, user)
	if err != nil {
		return nil, err
	}
	value, err := chClient.Query(apitxn.QueryRequest{ChaincodeID: chainCodeID,
		Fcn: "invoke", Args: args})
	if err != nil {
		return nil, err
	}
	return value, err
}

// QueryChannel the state of an asset on a given channel
func (setup *ClientConfig) QueryChannel(channel fab.Channel, chainCodeID string,
	args [][]byte, eventhub fab.EventHub) ([]byte, error) {

	ds, err := setup.SDK.DiscoveryProvider().NewDiscoveryService(channel.Name())
	if err != nil {
		return nil, err
	}
	ss, err := setup.SDK.SelectionProvider().NewSelectionService(channel.Name())
	if err != nil {
		return nil, err
	}
	chClient, err := chclient.NewChannelClient(setup.Client, channel,
		ds,
		ss,
		eventhub)
	if err != nil {
		return nil, err
	}

	value, err := chClient.Query(apitxn.QueryRequest{ChaincodeID: chainCodeID, Fcn: "invoke", Args: args})
	if err != nil {
		return nil, err
	}
	return value, err
}
