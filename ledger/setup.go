/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"log"
	"os"
	"path"

	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"

	deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"
	admin "github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/admin"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	Client          fab.FabricClient
	Channel         fab.Channel
	EventHub        fab.EventHub
	Sdk             deffab.FabricSDK
	ChannelClient   apitxn.ChannelClient
	ConnectEventHub bool
	ConfigFile      string
	OrgID           string
	OrdererOrgID    string
	ChannelID       string
	ChainCodeID     string
	Initialized     bool
	ChannelConfig   string
	UserID          string
	AdminUserID     string
	AdminUser       ca.User
}

// ExampleCC query and transaction arguments
var queryArgs = [][]byte{[]byte("query"), []byte("b")}
var txArgs = [][]byte{[]byte("move"), []byte("a"), []byte("b"), []byte("1")}

// ExampleCC init and upgrade args
var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
var upgradeArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("400")}

// ExampleCCQueryArgs returns example cc query args
func ExampleCCQueryArgs() [][]byte {
	return queryArgs
}

// ExampleCCTxArgs returns example cc move funds args
func ExampleCCTxArgs() [][]byte {
	return txArgs
}

//ExampleCCInitArgs returns example cc initialization args
func ExampleCCInitArgs() [][]byte {
	return initArgs
}

// Initialize reads configuration from file and sets up client, channel and event hub
func (setup *BaseSetupImpl) Initialize() error {
	// Create SDK setup for the integration tests
	sdkOptions := deffab.Options{
		ConfigFile: setup.ConfigFile,
	}

	sdk, err := deffab.NewSDK(sdkOptions)
	if err != nil {
		return errors.WithMessage(err, "SDK init failed")
	}

	setup.Sdk = *sdk

	session, err := sdk.NewPreEnrolledUserSession(setup.OrgID, setup.AdminUserID)
	if err != nil {
		return errors.WithMessage(err, "failed getting admin user session for org")
	}

	sc, err := sdk.NewSystemClient(session)
	if err != nil {
		return errors.WithMessage(err, "NewSystemClient failed")
	}

	setup.Client = sc
	setup.AdminUser = session.Identity()

	channel, err := setup.GetChannel(setup.Client, setup.ChannelID, []string{setup.OrgID})
	if err != nil {
		return errors.Wrapf(err, "create channel (%s) failed: %v", setup.ChannelID)
	}
	setup.Channel = channel

	chClient, err := sdk.NewChannelClient(setup.ChannelID, setup.UserID)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err)
	}

	setup.ChannelClient = chClient

	// Check if primary peer has joined channel
	alreadyJoined, err := HasPrimaryPeerJoinedChannel(sc, channel)
	if err != nil {
		return errors.WithMessage(err, "failed while checking if primary peer has already joined channel")
	}

	if !alreadyJoined {
		if err = admin.JoinChannel(sc, setup.AdminUser, channel); err != nil {
			return errors.WithMessage(err, "JoinChannel failed")
		}
	}

	if err := setup.setupEventHub(sc); err != nil {
		return err
	}

	setup.Initialized = true

	return nil
}

func (setup *BaseSetupImpl) setupEventHub(client fab.FabricClient) error {
	eventHub, err := setup.getEventHub(client)
	if err != nil {
		return err
	}

	if setup.ConnectEventHub {
		if err := eventHub.Connect(); err != nil {
			return errors.WithMessage(err, "eventHub connect failed")
		}
	}
	setup.EventHub = eventHub

	return nil
}

// InitConfig ...
func (setup *BaseSetupImpl) InitConfig() (apiconfig.Config, error) {
	configImpl, err := config.InitConfig(setup.ConfigFile)
	if err != nil {
		return nil, err
	}
	return configImpl, nil
}

// InstantiateCC ...
func (setup *BaseSetupImpl) InstantiateCC(chainCodeID string, chainCodePath string, chainCodeVersion string, args [][]byte) error {

	chaincodePolicy := cauthdsl.SignedByMspMember(setup.Client.UserContext().MspID())

	return admin.SendInstantiateCC(setup.Channel, chainCodeID, args, chainCodePath, chainCodeVersion, chaincodePolicy, []apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, setup.EventHub)
}

// UpgradeCC ...
func (setup *BaseSetupImpl) UpgradeCC(chainCodeID string, chainCodePath string, chainCodeVersion string, args [][]byte) error {

	chaincodePolicy := cauthdsl.SignedByMspMember(setup.Client.UserContext().MspID())

	return admin.SendUpgradeCC(setup.Channel, chainCodeID, args, chainCodePath, chainCodeVersion, chaincodePolicy, []apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, setup.EventHub)
}

// InstallCC ...
func (setup *BaseSetupImpl) InstallCC(chainCodeID string, chainCodePath string, chainCodeVersion string, chaincodePackage []byte) error {

	if err := admin.SendInstallCC(setup.Client, chainCodeID, chainCodePath, chainCodeVersion, chaincodePackage, peer.PeersToTxnProcessors(setup.Channel.Peers()), setup.GetDeployPath()); err != nil {
		return errors.WithMessage(err, "SendInstallProposal failed")
	}

	return nil
}

// GetDeployPath ..
func (setup *BaseSetupImpl) GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../fixtures/testdata")
}

// InstallAndInstantiateExampleCC ..
func (setup *BaseSetupImpl) InstallAndInstantiateExampleCC() error {

	chainCodePath := "github.com/example_cc"
	chainCodeVersion := "v0"

	if setup.ChainCodeID == "" {
		setup.ChainCodeID = GenerateRandomID()
	}

	if err := setup.InstallCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, nil); err != nil {
		return err
	}

	return setup.InstantiateCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, initArgs)
}

// UpgradeExampleCC ..
func (setup *BaseSetupImpl) UpgradeExampleCC() error {

	chainCodePath := "github.com/example_cc"
	chainCodeVersion := "v1"

	if setup.ChainCodeID == "" {
		setup.ChainCodeID = GenerateRandomID()
	}

	if err := setup.InstallCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, nil); err != nil {
		return err
	}

	return setup.UpgradeCC(setup.ChainCodeID, chainCodePath, chainCodeVersion, upgradeArgs)
}

// GetChannel initializes and returns a channel based on config
func (setup *BaseSetupImpl) GetChannel(client fab.FabricClient, channelID string, orgs []string) (fab.Channel, error) {

	channel, err := client.NewChannel(channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "NewChannel failed")
	}

	ordererConfig, err := client.Config().RandomOrdererConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "RandomOrdererConfig failed")
	}

	orderer, err := orderer.NewOrdererFromConfig(ordererConfig, client.Config())
	if err != nil {
		return nil, errors.WithMessage(err, "NewOrderer failed")
	}
	err = channel.AddOrderer(orderer)
	if err != nil {
		return nil, errors.WithMessage(err, "adding orderer failed")
	}

	for _, org := range orgs {
		peerConfig, err := client.Config().PeersConfig(org)
		if err != nil {
			return nil, errors.WithMessage(err, "reading peer config failed")
		}
		for _, p := range peerConfig {
			endorser, err := deffab.NewPeerFromConfig(&p, client.Config())
			if err != nil {
				return nil, errors.WithMessage(err, "NewPeer failed")
			}
			err = channel.AddPeer(endorser)
			if err != nil {
				return nil, errors.WithMessage(err, "adding peer failed")
			}
		}
	}

	return channel, nil
}

// CreateAndSendTransactionProposal ... TODO duplicate
func (setup *BaseSetupImpl) CreateAndSendTransactionProposal(channel fab.Channel, chainCodeID string,
	fcn string, args [][]byte, targets []apitxn.ProposalProcessor, transientData map[string][]byte) ([]*apitxn.TransactionProposalResponse, apitxn.TransactionID, error) {

	request := apitxn.ChaincodeInvokeRequest{
		Targets:      targets,
		Fcn:          fcn,
		Args:         args,
		TransientMap: transientData,
		ChaincodeID:  chainCodeID,
	}
	transactionProposalResponses, txnID, err := channel.SendTransactionProposal(request)
	if err != nil {
		return nil, txnID, err
	}

	for _, v := range transactionProposalResponses {
		if v.Err != nil {
			return nil, txnID, errors.Wrapf(v.Err, "endorser %s failed", v.Endorser)
		}
	}

	return transactionProposalResponses, txnID, nil
}

// CreateAndSendTransaction ...
func (setup *BaseSetupImpl) CreateAndSendTransaction(channel fab.Channel, resps []*apitxn.TransactionProposalResponse) (*apitxn.TransactionResponse, error) {

	tx, err := channel.CreateTransaction(resps)
	if err != nil {
		return nil, errors.WithMessage(err, "CreateTransaction failed")
	}

	transactionResponse, err := channel.SendTransaction(tx)
	if err != nil {
		return nil, errors.WithMessage(err, "SendTransaction failed")

	}

	if transactionResponse.Err != nil {
		return nil, errors.Wrapf(transactionResponse.Err, "orderer %s failed", transactionResponse.Orderer)
	}

	return transactionResponse, nil
}

// RegisterTxEvent registers on the given eventhub for the give transaction
// returns a boolean channel which receives true when the event is complete
// and an error channel for errors
// TODO - Duplicate
func (setup *BaseSetupImpl) RegisterTxEvent(txID apitxn.TransactionID, eventHub fab.EventHub) (chan bool, chan error) {
	done := make(chan bool)
	fail := make(chan error)

	eventHub.RegisterTxEvent(txID, func(txId string, errorCode pb.TxValidationCode, err error) {
		if err != nil {
			log.Printf("Received error event for txid(%s)", txId)
			fail <- err
		} else {
			log.Printf("Received success event for txid(%s)", txId)
			done <- true
		}
	})

	return done, fail
}

// getEventHub initilizes the event hub
func (setup *BaseSetupImpl) getEventHub(client fab.FabricClient) (fab.EventHub, error) {
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, errors.WithMessage(err, "NewEventHub failed")
	}
	foundEventHub := false
	peerConfig, err := client.Config().PeersConfig(setup.OrgID)
	if err != nil {
		return nil, errors.WithMessage(err, "PeersConfig failed")
	}
	for _, p := range peerConfig {
		if p.URL != "" {
			log.Printf("EventHub connect to peer (%s)", p.URL)
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
		return nil, errors.New("event hub configuration not found")
	}

	return eventHub, nil
}
