 import React from "react";
import { Flex, Text, Link } from "@chakra-ui/react";

 const Footer = (): React.ReactElement => {
  return (
    <Flex
      bg="system.gray"
      color="white"
      width="100%"
      justifyContent="center"
      alignItems="center"
      direction="column"
      padding={4}
    >
    
      <Flex width="100%" justifyContent="center" wrap="wrap">
         <Text width="100%" textAlign="center" color="white" fontSize="sm">
          A component of <Link href="https://trisa.io" color={"system.cyan"}>the TRISA architecture</Link> for Cryptocurrency Travel Rule compliance.
        </Text>
         <Text width="100%" textAlign="center" color="white" fontSize="sm">
Created and maintained by <Link href="https://rotational.io" color={"system.cyan"}> Rotational Labs</Link> in partnership with <Link href="https://cyphertrace.com" color={"system.cyan"}> CipherTrace</Link> CipherTrace on behalf of <Link href="https://trisa.io" color={"system.cyan"}>TRISA</Link> .
        </Text>
      </Flex>
    </Flex>
  );
 };

 export default Footer;
