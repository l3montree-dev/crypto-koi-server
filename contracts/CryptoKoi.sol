// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/token/ERC712/ERC712.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';

contract CryptoKoi is ERC712, AccessControl {
    constructor(string memory name, string memory symbol)
        ERC721(name, symbol)
    {
        // mark the creator of the contract as admin.
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function redeem(
        address account,
        uint256 tokenId,
        bytes calldata signature
    ) external {
        require(
            _verify(_hash(account, tokenId), signature),
            'Invalid signature'
        );
        _safeMint(account, tokenId);
    }

    function _baseURI()
        internal
        view
        virtual
        returns (string memory)
    {
        return 'https://api.crypto-koi.io/v1/tokens/';
    }

    function _hash(address account, uint256 tokenId)
        internal
        pure
        returns (bytes32)
    {
        return
            ECDSA.toEthSignedMessageHash(
                keccak256(abi.encodePacked(tokenId, account))
            );
    }

    function _verify(bytes32 digest, bytes memory signature)
        internal
        view
        returns (bool)
    {
        return
            hasRole(
                DEFAULT_ADMIN_ROLE,
                ECDSA.recover(digest, signature)
            );
    }
}
