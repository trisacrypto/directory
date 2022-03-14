import React, { FC } from "react";
import { Box, Text, Button } from "@chakra-ui/react";

const NeedsAttention = () => {
  return (
    <Box
      bg={"#D8EAF6"}
      minHeight={67}
      minWidth={246}
      pt={5}
      mt={10}
      mx={5}
      px={5}
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      fontSize={18}
    >
      <Box pb={2} display={"flex"} justifyContent={"space-between"}>
        <Text fontWeight={"bold"}> Needs Attention </Text>
        <Text> Complete Testnet Registration </Text>
        <Button
          bg={"#55ACD8"}
          color={"white"}
          height={34}
          width={142}
          _hover={{
            bg: "#55ACD8",
          }}
          _focus={{
            borderColor: "transparent",
          }}
        >
          Start
        </Button>
      </Box>
    </Box>
  );
};

export default NeedsAttention;
