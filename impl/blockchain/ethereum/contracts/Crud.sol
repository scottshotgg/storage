pragma solidity ^0.5.2;

// I don't think this works in Go... It fails to generate the ABI bindings...
pragma experimental ABIEncoderV2;

contract Crud {
    struct item {
        string id;
        bytes value;
        int64 timestamp;
        string[] keys;
        bool deleted;
    }

    mapping (bytes32 => item) itemMap;
    string[] itemKeys;


    function setItem(bytes32 domain, item memory obj) public returns (bool) {
        itemKeys.push(obj.id);
        itemMap[domain] = obj;
    }

    function getItem(bytes32 domain) public view returns (item memory) {
        return itemMap[domain];
    }

    function deleteItem(bytes32 domain) public {
        delete itemMap[domain];
    }
}