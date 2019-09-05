pragma solidity ^0.4.14;

/*
Example of how to verify BLS signatures and BGLS aggregate signatures in Ethereum.

Signatures are generated using https://github.com/Project-Arda/bgls
Code is based on https://github.com/Project-Arda/bgls-on-evm
*/

contract BLSE {
    struct G1Point {
        uint X;
        uint Y;
    }
    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint[2] X;
        uint[2] Y;
    }

    /// @return the generator of G1
    function P1() internal returns (G1Point) {
        return G1Point(1, 2);
    }

    /// @return the generator of G2
    function P2() internal returns (G2Point) {
        return G2Point(
            [11559732032986387107991004021392285783925812861821192530917403151452391805634,
            10857046999023057135944570762232829481370756359578518086990519993285655852781],

            [4082367875863433681332203403145435568316851327593401208105741076214120093531,
            8495653923123431417604973247489272438418190587263600148770280649306958101930]
        );
    }

    //Example of BGLS signature verification with 3 signers
    //Note that the messages differ in their last character.
    function verifyBGLS3(int i) public returns (bool) {

        uint numberOfSigners = uint(i);

        bytes memory message = hex"7b0a2020226f70656e223a207b0a20202020227072696365223a2039353931372c0a202020202274696d65223a207b0a20202020202022756e6978223a20313438333134323430302c0a2020202020202269736f223a2022323031362d31322d33315430303a30303a30302e3030305a220a202020207d0a20207d2c0a202022636c6f7365223a207b0a20202020227072696365223a2039363736302c0a202020202274696d65223a207b0a20202020202022756e6978223a20313438333232383830302c0a2020202020202269736f223a2022323031372d30312d30315430303a30303a30302e3030305a220a202020207d0a20207d2c0a2020226c6f6f6b7570223a207b0a20202020227072696365223a2039363736302c0a20202020226b223a20312c0a202020202274696d65223a207b0a20202020202022756e6978223a20313438333232383830302c0a2020202020202269736f223a2022323031372d30312d30315430303a30303a30302e3030305a220a202020207d0a20207d0a7d0a6578616d706c652e636f6d2f6170692f31";

        G1Point memory signature = G1Point(385846518441062319503502284295243290270560187383398932887791670182362540842, 19731933537428695151702009864745685458233056709189425720845387511061953267292);

        G2Point memory v = G2Point(
            [1787282038370667094324364195810339512415273589223814213215040505578200405366, 414568866548933554513940840943382696902163788831396286279770126458218272940],
            [6560020551439455112781785895092032589010633560844445112872109862153018855017, 19411093226570397520343120724285433000937737461010544490862811136406407315543]
        );

        G1Point memory h = hashToG1(message);

        G1Point[] memory a = new G1Point[](numberOfSigners + 1);
        G2Point[] memory b = new G2Point[](numberOfSigners + 1);

        a[0] = negate(signature);
        b[0] = P2();

        for (uint index = 1; index <= numberOfSigners; index ++) {
            a[index] = h;
            b[index] = v;
        }

        return pairing(a, b);
    }
    
    function testAdd(uint number) public returns(uint) {
        uint sum = 0;
        for (uint i; i < number; i++) {
            sum += i;
        }
        
        return sum;
    }

    /// @return the result of computing the pairing check
    /// e(p1[0], p2[0]) *  .... * e(p1[n], p2[n]) == 1
    /// For example pairing([P1(), P1().negate()], [P2(), P2()]) should
    /// return true.
    function pairing(G1Point[] p1, G2Point[] p2) internal returns (bool) {
        require(p1.length == p2.length);
        uint elements = p1.length;
        uint inputSize = elements * 6;
        uint[] memory input = new uint[](inputSize);

        for (uint i = 0; i < elements; i++)
        {
            input[i * 6 + 0] = p1[i].X;
            input[i * 6 + 1] = p1[i].Y;
            input[i * 6 + 2] = p2[i].X[0];
            input[i * 6 + 3] = p2[i].X[1];
            input[i * 6 + 4] = p2[i].Y[0];
            input[i * 6 + 5] = p2[i].Y[1];
        }

        uint[1] memory out;
        bool success;

        assembly {
            success := call(sub(gas, 2000), 8, 0, add(input, 0x20), mul(inputSize, 0x20), out, 0x20)
        // Use "invalid" to make gas estimation work
            switch success case 0 {invalid}
        }
        require(success);
        return out[0] != 0;
    }

    /// Convenience method for a pairing check for two pairs.
    function pairing2(G1Point a1, G2Point a2, G1Point b1, G2Point b2) internal returns (bool) {
        G1Point[] memory p1 = new G1Point[](2);
        G2Point[] memory p2 = new G2Point[](2);
        p1[0] = a1;
        p1[1] = b1;
        p2[0] = a2;
        p2[1] = b2;
        return pairing(p1, p2);
    }

    function hashToG1(bytes message) internal returns (G1Point) {
        uint256 h = uint256(keccak256(message));
        return mul(P1(), h);
    }

    /// @return the negation of p, i.e. p.add(p.negate()) should be zero.
    function negate(G1Point p) internal returns (G1Point) {
        // The prime q in the base field F_q for G1
        uint q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
        if (p.X == 0 && p.Y == 0)
            return G1Point(0, 0);
        return G1Point(p.X, q - (p.Y % q));
    }

    /// @return the sum of two points of G1
    function add(G1Point p1, G1Point p2) internal returns (G1Point r) {
        uint[4] memory input;
        input[0] = p1.X;
        input[1] = p1.Y;
        input[2] = p2.X;
        input[3] = p2.Y;
        bool success;
        assembly {
            success := call(sub(gas, 2000), 6, 0, input, 0xc0, r, 0x60)
        // Use "invalid" to make gas estimation work
            switch success case 0 {invalid}
        }
        require(success);
    }
    /// @return the product of a point on G1 and a scalar, i.e.
    /// p == p.mul(1) and p.add(p) == p.mul(2) for all points p.
    function mul(G1Point p, uint s) internal returns (G1Point r) {
        uint[3] memory input;
        input[0] = p.X;
        input[1] = p.Y;
        input[2] = s;
        bool success;
        assembly {
            success := call(sub(gas, 2000), 7, 0, input, 0x80, r, 0x60)
        // Use "invalid" to make gas estimation work
            switch success case 0 {invalid}
        }
        require(success);
    }
}