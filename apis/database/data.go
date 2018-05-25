package database

import (
	"encoding/json"
	"github.com/scorum/openledger-go/types"
)

type Asset struct {
	ID                 types.ObjectID `json:"id"`
	Symbol             string         `json:"symbol"`
	Precision          uint8          `json:"precision"`
	Issuer             string         `json:"issuer"`
	DynamicAssetDataID string         `json:"dynamic_asset_data_id"`
}

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             types.Time        `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
}

type Block struct {
	TransactionMerkleRoot string              `json:"transaction_merkle_root"`
	Previous              string              `json:"previous"`
	Timestamp             types.Time          `json:"timestamp"`
	Witness               string              `json:"witness"`
	Extensions            []json.RawMessage   `json:"extensions"`
	WitnessSignature      string              `json:"witness_signature"`
	Transactions          []types.Transaction `json:"transactions"`
}

type MarketTicker struct {
	Time          types.Time     `json:"time"`
	Base          types.ObjectID `json:"base"`
	Quote         types.ObjectID `json:"quote"`
	Latest        string         `json:"latest"`
	LowestAsk     string         `json:"lowest_ask"`
	HighestBid    string         `json:"highest_bid"`
	PercentChange string         `json:"percent_change"`
	BaseVolume    string         `json:"base_volume"`
	QuoteVolume   string         `json:"quote_volume"`
}

type LimitOrder struct {
	ID          types.ObjectID `json:"id"`
	Expiration  types.Time     `json:"expiration"`
	Seller      types.ObjectID `json:"seller"`
	ForSale     uint64         `json:"for_sale"`
	DeferredFee uint64         `json:"deferred_fee"`
	SellPrice   types.Price    `json:"sell_price"`
}

type DynamicGlobalProperties struct {
	ID                             types.ObjectID `json:"id"`
	HeadBlockNumber                uint32         `json:"head_block_number"`
	HeadBlockID                    string         `json:"head_block_id"`
	Time                           types.Time     `json:"time"`
	CurrentWitness                 types.ObjectID `json:"current_witness"`
	NextMaintenanceTime            types.Time     `json:"next_maintenance_time"`
	LastBudgetTime                 types.Time     `json:"last_budget_time"`
	AccountsRegisteredThisInterval int            `json:"accounts_registered_this_interval"`
	DynamicFlags                   int            `json:"dynamic_flags"`
	RecentSlotsFilled              string         `json:"recent_slots_filled"`
	LastIrreversibleBlockNum       uint32         `json:"last_irreversible_block_num"`
	CurrentAslot                   int64          `json:"current_aslot"`
	WitnessBudget                  int64          `json:"witness_budget"`
	RecentlyMissedCount            int64          `json:"recently_missed_count"`
}

type Config struct {
	GrapheneSymbol               string `json:"GRAPHENE_SYMBOL"`
	GrapheneAddressPrefix        string `json:"GRAPHENE_ADDRESS_PREFIX"`
	GrapheneMinAccountNameLength uint8  `json:"GRAPHENE_MIN_ACCOUNT_NAME_LENGTH"`
	GrapheneMaxAccountNameLength uint8  `json:"GRAPHENE_MAX_ACCOUNT_NAME_LENGTH"`
	GrapheneMinAssetSymbolLength uint8  `json:"GRAPHENE_MIN_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxAssetSymbolLength uint8  `json:"GRAPHENE_MAX_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxShareSupply       string `json:"GRAPHENE_MAX_SHARE_SUPPLY"`
}

type AccountsMap map[string]types.ObjectID

func (o *AccountsMap) UnmarshalJSON(b []byte) error {
	out := make(map[string]types.ObjectID)

	// unmarshal array
	var arr []json.RawMessage
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}

	var (
		key string
		obj types.ObjectID
	)

	for _, item := range arr {
		account := []interface{}{&key, &obj}
		if err := json.Unmarshal(item, &account); err != nil {
			return err
		}

		out[key] = obj
	}

	*o = out
	return nil
}

/*
{
  "GRAPHENE_MAX_PAY_RATE": 10000,
  "GRAPHENE_MAX_SIG_CHECK_DEPTH": 2,
  "GRAPHENE_MIN_TRANSACTION_SIZE_LIMIT": 1024,
  "GRAPHENE_MIN_BLOCK_INTERVAL": 1,
  "GRAPHENE_MAX_BLOCK_INTERVAL": 30,
  "GRAPHENE_DEFAULT_BLOCK_INTERVAL": 5,
  "GRAPHENE_DEFAULT_MAX_TRANSACTION_SIZE": 2048,
  "GRAPHENE_DEFAULT_MAX_BLOCK_SIZE": 2000000,
  "GRAPHENE_DEFAULT_MAX_TIME_UNTIL_EXPIRATION": 86400,
  "GRAPHENE_DEFAULT_MAINTENANCE_INTERVAL": 86400,
  "GRAPHENE_DEFAULT_MAINTENANCE_SKIP_SLOTS": 3,
  "GRAPHENE_MIN_UNDO_HISTORY": 10,
  "GRAPHENE_MAX_UNDO_HISTORY": 10000,
  "GRAPHENE_MIN_BLOCK_SIZE_LIMIT": 5120,
  "GRAPHENE_MIN_TRANSACTION_EXPIRATION_LIMIT": 150,
  "GRAPHENE_BLOCKCHAIN_PRECISION": 100000,
  "GRAPHENE_BLOCKCHAIN_PRECISION_DIGITS": 5,
  "GRAPHENE_DEFAULT_TRANSFER_FEE": 100000,
  "GRAPHENE_MAX_INSTANCE_ID": "281474976710655",
  "GRAPHENE_100_PERCENT": 10000,
  "GRAPHENE_1_PERCENT": 100,
  "GRAPHENE_MAX_MARKET_FEE_PERCENT": 10000,
  "GRAPHENE_DEFAULT_FORCE_SETTLEMENT_DELAY": 86400,
  "GRAPHENE_DEFAULT_FORCE_SETTLEMENT_OFFSET": 0,
  "GRAPHENE_DEFAULT_FORCE_SETTLEMENT_MAX_VOLUME": 2000,
  "GRAPHENE_DEFAULT_PRICE_FEED_LIFETIME": 86400,
  "GRAPHENE_MAX_FEED_PRODUCERS": 200,
  "GRAPHENE_DEFAULT_MAX_AUTHORITY_MEMBERSHIP": 10,
  "GRAPHENE_DEFAULT_MAX_ASSET_WHITELIST_AUTHORITIES": 10,
  "GRAPHENE_DEFAULT_MAX_ASSET_FEED_PUBLISHERS": 10,
  "GRAPHENE_COLLATERAL_RATIO_DENOM": 1000,
  "GRAPHENE_MIN_COLLATERAL_RATIO": 1001,
  "GRAPHENE_MAX_COLLATERAL_RATIO": 32000,
  "GRAPHENE_DEFAULT_MAINTENANCE_COLLATERAL_RATIO": 1750,
  "GRAPHENE_DEFAULT_MAX_SHORT_SQUEEZE_RATIO": 1500,
  "GRAPHENE_DEFAULT_MARGIN_PERIOD_SEC": 2592000,
  "GRAPHENE_DEFAULT_MAX_WITNESSES": 1001,
  "GRAPHENE_DEFAULT_MAX_COMMITTEE": 1001,
  "GRAPHENE_DEFAULT_MAX_PROPOSAL_LIFETIME_SEC": 2419200,
  "GRAPHENE_DEFAULT_COMMITTEE_PROPOSAL_REVIEW_PERIOD_SEC": 1209600,
  "GRAPHENE_DEFAULT_NETWORK_PERCENT_OF_FEE": 2000,
  "GRAPHENE_DEFAULT_LIFETIME_REFERRER_PERCENT_OF_FEE": 3000,
  "GRAPHENE_DEFAULT_MAX_BULK_DISCOUNT_PERCENT": 5000,
  "GRAPHENE_DEFAULT_BULK_DISCOUNT_THRESHOLD_MIN": 100000000,
  "GRAPHENE_DEFAULT_BULK_DISCOUNT_THRESHOLD_MAX": "10000000000",
  "GRAPHENE_DEFAULT_CASHBACK_VESTING_PERIOD_SEC": 31536000,
  "GRAPHENE_DEFAULT_CASHBACK_VESTING_THRESHOLD": 10000000,
  "GRAPHENE_DEFAULT_BURN_PERCENT_OF_FEE": 2000,
  "GRAPHENE_WITNESS_PAY_PERCENT_PRECISION": 1000000000,
  "GRAPHENE_DEFAULT_MAX_ASSERT_OPCODE": 1,
  "GRAPHENE_DEFAULT_FEE_LIQUIDATION_THRESHOLD": 10000000,
  "GRAPHENE_DEFAULT_ACCOUNTS_PER_FEE_SCALE": 1000,
  "GRAPHENE_DEFAULT_ACCOUNT_FEE_SCALE_BITSHIFTS": 4,
  "GRAPHENE_MAX_WORKER_NAME_LENGTH": 63,
  "GRAPHENE_MAX_URL_LENGTH": 127,
  "GRAPHENE_NEAR_SCHEDULE_CTR_IV": "7640891576956012808",
  "GRAPHENE_FAR_SCHEDULE_CTR_IV": "13503953896175478587",
  "GRAPHENE_CORE_ASSET_CYCLE_RATE": 17,
  "GRAPHENE_CORE_ASSET_CYCLE_RATE_BITS": 32,
  "GRAPHENE_DEFAULT_WITNESS_PAY_PER_BLOCK": 1000000,
  "GRAPHENE_DEFAULT_WITNESS_PAY_VESTING_SECONDS": 86400,
  "GRAPHENE_DEFAULT_WORKER_BUDGET_PER_DAY": "50000000000",
  "GRAPHENE_MAX_INTEREST_APR": 10000,
  "GRAPHENE_COMMITTEE_ACCOUNT": "1.2.0",
  "GRAPHENE_WITNESS_ACCOUNT": "1.2.1",
  "GRAPHENE_RELAXED_COMMITTEE_ACCOUNT": "1.2.2",
  "GRAPHENE_NULL_ACCOUNT": "1.2.3",
  "GRAPHENE_TEMP_ACCOUNT": "1.2.4"
}
*/
