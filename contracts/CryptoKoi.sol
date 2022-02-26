// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/token/ERC721/ERC721.sol';
import '@openzeppelin/contracts/utils/cryptography/ECDSA.sol';

contract CryptoKoi is ERC721 {
    address payable public owner;

    string baseURI;
    uint256 price;

    constructor(
        string memory name,
        string memory symbol,
        string memory uri,
        uint256 p
    ) ERC721(name, symbol) {
        owner = payable(msg.sender);
        baseURI = uri;
        price = p;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, 'Owner privilege only');
        _;
    }

    function withdrawAll() external onlyOwner {
        uint256 amount = address(this).balance;

        (bool success, ) = owner.call{value: amount}('');

        require(success, 'withdrawAll: Transfer failed');

        emit Transfer(address(0), owner, amount);
    }

    function killSwitch() external onlyOwner {
        selfdestruct(owner);
    }

    function setPrice(uint256 p) external onlyOwner {
        price = p;
    }

    function getPrice() external view returns (uint256) {
        return price;
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(ERC721)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function redeem(
        address account,
        uint256 tokenId,
        bytes calldata signature
    ) external payable {
        require(msg.value >= price, 'Insufficient funds');
        require(
            _verify(_hash(account, tokenId), signature),
            'Invalid signature'
        );

        _safeMint(account, tokenId);
    }

    function setBaseURI(string calldata uri) external onlyOwner {
        baseURI = uri;
    }

    function _baseURI()
        internal
        view
        virtual
        override(ERC721)
        returns (string memory)
    {
        return baseURI;
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
        return owner == ECDSA.recover(digest, signature);
    }
}
