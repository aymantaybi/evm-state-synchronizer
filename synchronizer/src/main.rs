use alloy::{
    eips::BlockNumberOrTag,
    primitives::address,
    providers::{ext::DebugApi, ProviderBuilder},
    rpc::types::{
        trace::geth::{
            GethDebugBuiltInTracerType, GethDebugTracerType, GethDebugTracingCallOptions,
        },
        TransactionRequest,
    },
    sol,
    sol_types::SolCall,
};
use IERC20::balanceOfCall;

sol!(
    #[allow(missing_docs)]
    #[sol(rpc)]
    #[derive(Debug)]
    IERC20,
    "src/abi/IERC20.json"
);

#[tokio::main]
async fn main() {
    let account_address = address!("c1eb47de5d549d45a871e32d9d082e7ac5d2e3ed");
    let token_address = address!("0b7007c13325c48911f73a2dad5fa5dcbf808adc");

    let rpc_url =
        "https://lb.drpc.org/ogrpc?network=ronin&dkey=Aln7W35A6Edpse2U-unSi3ffsd8RUugR77V9vmJKmvm9";

    let provider = ProviderBuilder::new().on_http(rpc_url.parse().unwrap());

    let mut tx = TransactionRequest::default()
        .to(token_address)
        .from(account_address);

    let encoded_balance_of_call = balanceOfCall {
        _owner: account_address.clone(),
    }
    .abi_encode();

    tx.input.data = Some(encoded_balance_of_call.into());

    let mut trace_options = GethDebugTracingCallOptions::default();
    trace_options.tracing_options.tracer = Some(GethDebugTracerType::BuiltInTracer(
        GethDebugBuiltInTracerType::PreStateTracer,
    ));

    let result = provider
        .debug_trace_call(tx, BlockNumberOrTag::Latest, trace_options)
        .await;

    println!("{:#?}", result);
}
