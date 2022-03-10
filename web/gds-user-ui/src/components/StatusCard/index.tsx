import React from "react";
import { Stack, Box, Text, Heading, HStack } from "@chakra-ui/react";

interface StatusCardProps {
  testnetstatus: string;
  mainnetstatus: string;
}

const StatusCard = ({ testnetstatus, mainnetstatus }: StatusCardProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      color={"#252733"}
      height={167}
      maxWidth={451}
      fontSize={18}
      p={5}
      mt={10}
      px={5}
    >
      <Stack>
        <Heading fontSize={20}>Certification Status</Heading>
        <HStack spacing={10}>
          <Text>Testnet</Text>
          <Text>{testnetstatus}</Text>
        </HStack>
        <HStack spacing={8}>
          <Text>Mainnet</Text>
          <Text>{mainnetstatus}</Text>
        </HStack>
      </Stack>
    </Box>
  );
};
StatusCard.defaultProps = {
  testnetstatus: "In progress",
  mainnetstatus: "Not Eligible yet ",
};

export default StatusCard;
