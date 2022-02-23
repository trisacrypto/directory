import React from "react";
import { Flex, Text, HStack } from "@chakra-ui/react";

 const LandingFooter = (): React.ReactElement => {
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
      <HStack spacing={8} mb={8}>
        <Text width="100%" textAlign="center" color="white" fontSize="sm">
         A component of the TRISA architecture for Cryptocurrency Travel Rule compliance.
Created and maintained by Rotational Labs in partnership with CipherTrace on behalf of TRISA.
        </Text>
      </HStack>
      <Flex width="100%" justifyContent="center" wrap="wrap">
       
      </Flex>
    </Flex>
  );
 };

 export default LandingFooter;
