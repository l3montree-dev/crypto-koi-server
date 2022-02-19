import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { expect } from 'chai';
import { Contract, Signer, Wallet } from 'ethers';
import { ethers } from 'hardhat';
import tokens from './tokens.json';

const otherUserAddress = '0xa111C225A0aFd5aD64221B1bc1D5d817e5D3Ca15';
const privateKey =
  '0xc0c1e7d82fae79ce7727bd94e3e74deafbce52fc5618d9fd5557f41e83d4c149';
const tokenId = '239264596381739575473221873891232270519';
const expectedHexHash =
  '8e130016394d7e04194944001ec36c64de13a73e81c716cd10016d5d83347a00';

const expectedSignature =
  '0x0577530589f065fdb25b8f29132865782ab2a4ea75a294ba56deecddeeefb77b18755f1811bb76dfadf417ff58f6bd2b593ddb4c80b1eaa85752e0df5a5b44f41b';

type TokenType = {
  tokenId?: string;
  account?: string;
  signature?: string;
};

async function deploy(
  name: string,
  signer: Signer,
  ...params: unknown[]
) {
  const Contract = await ethers.getContractFactory(name, signer);
  return await Contract.deploy(...params).then((f) => f.deployed());
}

function hashToken(tokenId: string, account: string) {
  return Buffer.from(
    ethers.utils
      .solidityKeccak256(['uint256', 'address'], [tokenId, account])
      .slice(2),
    'hex',
  );
}

describe('CryptoKoi', function () {
  let accounts: SignerWithAddress[];
  let signer: SignerWithAddress;
  before(async function () {
    accounts = await ethers.getSigners();
    signer = accounts[0];
  });

  describe('Golang integration', () => {
    let contract: Contract;
    let admin: Wallet;
    before(async () => {
      contract = await deploy('CryptoKoi', admin, 'Name', 'Symbol');

      admin = new ethers.Wallet(privateKey);
    });

    it('should work with fixed values', async () => {
      // the fixed values are provided by the golang executable. The values can be found: web3_test.go
      const hash = hashToken(tokenId, otherUserAddress);

      expect(hash.toString('hex')).to.eq(expectedHexHash);

      const signature = await admin.signMessage(hash);
      expect(signature).to.eq(expectedSignature);

      // check if it is possible to redeem the token now with the provided values.
      // fist deploy the smart contract using the admin account.

      await expect(
        contract
          .connect(accounts[2])
          .redeem(otherUserAddress, tokenId, signature),
      )
        .to.emit(contract, 'Transfer')
        .withArgs(
          ethers.constants.AddressZero,
          otherUserAddress,
          tokenId,
        );
    });
  });

  describe('Mint all elements', function () {
    let contract: Contract;

    before(async function () {
      contract = await deploy('CryptoKoi', signer, 'Name', 'Symbol');
    });

    for (const [tokenId, account] of Object.entries(tokens)) {
      it(
        'Element: ' + tokenId + ' should be minted to ' + account,
        async function () {
          /**
           * Account[0] (minter) creates signature
           */
          const signature = await signer.signMessage(
            hashToken(tokenId, account),
          );
          /**
           * Account[2] (anyone?) redeems token using signature
           */
          await expect(
            contract
              .connect(accounts[2])
              .redeem(account, tokenId, signature),
          )
            .to.emit(contract, 'Transfer')
            .withArgs(ethers.constants.AddressZero, account, tokenId);
        },
      );
    }
  });

  describe('Duplicate mint', function () {
    let contract: Contract;
    let token: TokenType;
    before(async function () {
      contract = await deploy('CryptoKoi', signer, 'Name', 'Symbol');

      token = {};
      const t = Object.entries(tokens).find(Boolean);
      if (t) {
        [token.tokenId, token.account] = t;
        token.signature = await signer.signMessage(
          hashToken(token.tokenId, token.account),
        );
      }
    });

    it('mint once - success', async function () {
      await expect(
        contract.redeem(
          token.account,
          token.tokenId,
          token.signature,
        ),
      )
        .to.emit(contract, 'Transfer')
        .withArgs(
          ethers.constants.AddressZero,
          token.account,
          token.tokenId,
        );
    });

    it('should return the correct current balance', async () => {
      await contract.redeem(
        token.account,
        token.tokenId,
        token.signature,
      );

      expect(await contract.balanceOf(token.account)).to.equal(
        BigInt(1),
      );
    });

    it('should return the correct owner', async () => {
      await contract.redeem(
        token.account,
        token.tokenId,
        token.signature,
      );

      expect(await contract.ownerOf(token.tokenId)).to.equal(
        token.account,
      );
    });

    it('mint twice - failure', async function () {
      await expect(
        contract.redeem(
          token.account,
          token.tokenId,
          token.signature,
        ),
      ).to.be.revertedWith('ERC721: token already minted');
    });
  });

  describe('Frontrun', function () {
    let contract: Contract;
    let token: TokenType;
    before(async function () {
      contract = await deploy('CryptoKoi', signer, 'Name', 'Symbol');

      token = {};
      const t = Object.entries(tokens).find(Boolean);
      if (t) {
        [token.tokenId, token.account] = t;
        token.signature = await accounts[1].signMessage(
          hashToken(token.tokenId, token.account),
        );
      }
    });

    it('Change owner - success', async function () {
      await expect(
        contract.redeem(
          accounts[0].address,
          token.tokenId,
          token.signature,
        ),
      ).to.be.revertedWith('Invalid signature');
    });
  });
});
