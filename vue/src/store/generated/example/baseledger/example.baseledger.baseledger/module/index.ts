// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgUpdateBaseledgerTransaction } from "./types/baseledger/tx";
import { MsgDeleteBaseledgerTransaction } from "./types/baseledger/tx";
import { MsgCreateBaseledgerTransaction } from "./types/baseledger/tx";


const types = [
  ["/example.baseledger.baseledger.MsgUpdateBaseledgerTransaction", MsgUpdateBaseledgerTransaction],
  ["/example.baseledger.baseledger.MsgDeleteBaseledgerTransaction", MsgDeleteBaseledgerTransaction],
  ["/example.baseledger.baseledger.MsgCreateBaseledgerTransaction", MsgCreateBaseledgerTransaction],
  
];

const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw new Error("wallet is required");

  const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee=defaultFee, memo=null }: SignAndBroadcastOptions) => memo?client.signAndBroadcast(address, msgs, fee,memo):client.signAndBroadcast(address, msgs, fee),
    msgUpdateBaseledgerTransaction: (data: MsgUpdateBaseledgerTransaction): EncodeObject => ({ typeUrl: "/example.baseledger.baseledger.MsgUpdateBaseledgerTransaction", value: data }),
    msgDeleteBaseledgerTransaction: (data: MsgDeleteBaseledgerTransaction): EncodeObject => ({ typeUrl: "/example.baseledger.baseledger.MsgDeleteBaseledgerTransaction", value: data }),
    msgCreateBaseledgerTransaction: (data: MsgCreateBaseledgerTransaction): EncodeObject => ({ typeUrl: "/example.baseledger.baseledger.MsgCreateBaseledgerTransaction", value: data }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
