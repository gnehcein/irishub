package main

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authutil "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	slashingcmd "github.com/cosmos/cosmos-sdk/x/slashing/client/cli"
	stakecmd "github.com/cosmos/cosmos-sdk/x/stake/client/cli"
	"github.com/irisnet/irishub/app"

	authcmd "github.com/irisnet/irishub/client/auth/cli"
	bankcmd "github.com/irisnet/irishub/client/bank/cli"
	govcmd "github.com/irisnet/irishub/client/gov/cli"
	upgradecmd "github.com/irisnet/irishub/client/upgrade/cli"
	"github.com/irisnet/irishub/version"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "iriscli",
		Short: "irishub command line interface",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// add standard rpc commands
	rpc.AddCommands(rootCmd)

	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tendermintCmd, cdc)

	rootCmd.AddCommand(
		tendermintCmd,
		client.LineBreak,
	)

	//Add stake commands
	stakeCmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake and validation subcommands",
	}
	stakeCmd.AddCommand(
		client.GetCommands(
			stakecmd.GetCmdQueryValidator("stake", cdc),
			stakecmd.GetCmdQueryValidators("stake", cdc),
			stakecmd.GetCmdQueryDelegation("stake", cdc),
			stakecmd.GetCmdQueryDelegations("stake", cdc),
			slashingcmd.GetCmdQuerySigningInfo("slashing", cdc),
		)...)
	stakeCmd.AddCommand(
		client.PostCommands(
			stakecmd.GetCmdCreateValidator(cdc),
			stakecmd.GetCmdEditValidator(cdc),
			stakecmd.GetCmdDelegate(cdc),
			stakecmd.GetCmdUnbond("stake", cdc),
			stakecmd.GetCmdRedelegate("stake", cdc),
			slashingcmd.GetCmdUnrevoke(cdc),
		)...)
	rootCmd.AddCommand(
		stakeCmd,
	)

	//Add gov commands
	govCmd := &cobra.Command{
		Use:   "gov",
		Short: "Governance and voting subcommands",
	}
	govCmd.AddCommand(
		client.GetCommands(
			govcmd.GetCmdQueryProposal("gov", cdc),
			govcmd.GetCmdQueryProposals("gov", cdc),
			govcmd.GetCmdQueryVote("gov", cdc),
			govcmd.GetCmdQueryVotes("gov", cdc),
			govcmd.GetCmdQueryConfig("iparams", cdc),
		)...)
	govCmd.AddCommand(
		client.PostCommands(
			govcmd.GetCmdSubmitProposal(cdc),
			govcmd.GetCmdDeposit(cdc),
			govcmd.GetCmdVote(cdc),
		)...)
	rootCmd.AddCommand(
		govCmd,
	)

	//Add upgrade commands
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Software Upgrade subcommands",
	}
	upgradeCmd.AddCommand(
		client.GetCommands(
			upgradecmd.GetCmdQuerySwitch("upgrade", cdc),
			upgradecmd.GetCmdInfo("upgrade", cdc),
		)...)
	upgradeCmd.AddCommand(
		client.PostCommands(
			upgradecmd.GetCmdSubmitSwitch(cdc),
		)...)
	rootCmd.AddCommand(
		upgradeCmd,
	)

	//Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authutil.GetAccountDecoder(cdc)),
			bankcmd.GetCmdQueryCoinType(cdc),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
	)
	rootCmd.AddCommand(
		client.GetCommands(
			version.GetCmdVersion("upgrade", cdc),
		)...)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
