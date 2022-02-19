/**
 * @type import('hardhat/config').HardhatUserConfig
 */
import '@nomiclabs/hardhat-waffle';
import '@nomiclabs/hardhat-ethers';
import 'dotenv/config';
import 'hardhat-gas-reporter';

const { HARDHAT_PORT, ALCHEMY_API_KEY } = process.env;

export = {
  gasReporter: {
    enabled: true,
    currency: 'EUR',
    token: 'MATIC',
    coinmarketcap: '0a636342-83a5-453d-837e-797d24235436',
  },
  solidity: '0.8.1',
  networks: {
    ropsten: {
      url: 'https://eth-ropsten.alchemyapi.io/v2/' + ALCHEMY_API_KEY,
    },
    localhost: { url: `http://127.0.0.1:${HARDHAT_PORT}` },
    hardhat: {
      chainId: 1337,
      accounts: [
        {
          privateKey:
            '0xc0c1e7d82fae79ce7727bd94e3e74deafbce52fc5618d9fd5557f41e83d4c149',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x998c1abd3a4ff680a14076000a442eb9f7c00ccc20513aee07e1d2466b2fe7ae',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0xc3365a2a4a29d54d99e337e24da30de2393d833b922bf903f807170b3773109f',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0xbe84e0703be0713b050a950f0d16494be9bafda841b4f4c251643ce537dd537a',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x4fae742b3d198c6b3378c069e71cdb8908c39740d51817b2a095991b7d1cab02',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x59e455c91ffa09d065b5a77db1c959bbfcddeb954c1f8253bfb8861aa3999544',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x0e9014cdb233d74feb7818cfae21f642abcea7d843b96e24ecd1b6748b2202a2',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0xd41a94b8753a8a310af57b8d11dff5fa351cf8d4516337c42d19c8bc8823542e',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x038000b9220ca2ff7f2dda5e4d65a094eb1c67a83d17fc86dda4bfb434dd2cd2',
          balance: '1000000000000000000000',
        },
        {
          privateKey:
            '0x33a9facefe26fcee475184c9c0f83bf0b5a692ff8dbf6f6094a51e7e524f1ed0',
          balance: '1000000000000000000000',
        },
      ],
    },
  },
  paths: {
    sources: './contracts',
    tests: './test',
    cache: './cache',
    artifacts: './artifacts',
  },
};
