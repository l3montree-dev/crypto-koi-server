import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers';
import { Contract } from 'ethers';
import { ethers } from 'hardhat';
import { expect } from 'chai';
import tokens from './tokens.json';

type TokenType = {
  tokenId?: string;
  account?: string;
  signature?: string;
};

async function deploy(name: string, ...params: unknown[]) {
  const Contract = await ethers.getContractFactory(name);
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
  before(async function () {
    accounts = await ethers.getSigners();
  });

  describe('Mint all elements', function () {
    let registry: Contract;

    before(async function () {
      registry = await deploy('CryptoKoi', 'Name', 'Symbol');
      await registry.grantRole(
        await registry.MINTER_ROLE(),
        accounts[1].address,
      );
    });

    for (const [tokenId, account] of Object.entries(tokens)) {
      it('element', async function () {
        /**
         * Account[1] (minter) creates signature
         */
        const signature = await accounts[1].signMessage(
          hashToken(tokenId, account),
        );
        /**
         * Account[2] (anyone?) redeems token using signature
         */
        await expect(
          registry
            .connect(accounts[2])
            .redeem(account, tokenId, signature),
        )
          .to.emit(registry, 'Transfer')
          .withArgs(ethers.constants.AddressZero, account, tokenId);
      });
    }
  });

  describe('Duplicate mint', function () {
    let registry: Contract;
    let token: TokenType;
    before(async function () {
      registry = await deploy('CryptoKoi', 'Name', 'Symbol');
      await registry.grantRole(
        await registry.MINTER_ROLE(),
        accounts[1].address,
      );

      token = {};
      const t = Object.entries(tokens).find(Boolean);
      if (t) {
        [token.tokenId, token.account] = t;
        token.signature = await accounts[1].signMessage(
          hashToken(token.tokenId, token.account),
        );
      }
    });

    it('mint once - success', async function () {
      await expect(
        registry.redeem(
          token.account,
          token.tokenId,
          token.signature,
        ),
      )
        .to.emit(registry, 'Transfer')
        .withArgs(
          ethers.constants.AddressZero,
          token.account,
          token.tokenId,
        );
    });

    it('mint twice - failure', async function () {
      await expect(
        registry.redeem(
          token.account,
          token.tokenId,
          token.signature,
        ),
      ).to.be.revertedWith('ERC721: token already minted');
    });
  });

  describe('Frontrun', function () {
    let registry: Contract;
    let token: TokenType;
    before(async function () {
      registry = await deploy('CryptoKoi', 'Name', 'Symbol');
      await registry.grantRole(
        await registry.MINTER_ROLE(),
        accounts[1].address,
      );

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
        registry.redeem(
          accounts[0].address,
          token.tokenId,
          token.signature,
        ),
      ).to.be.revertedWith('Invalid signature');
    });
  });
});
