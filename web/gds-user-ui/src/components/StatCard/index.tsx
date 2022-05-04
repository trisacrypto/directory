import React from "react";
import { Stack, Box, Text, Heading } from "@chakra-ui/react";

interface StatCardProps {
  title: string;
  number: number;
}
const StatCard = ({ title, number }: StatCardProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      color={"#252733"}
      height={167}
      maxWidth={246}
      fontSize={18}
      p={5}
      mt={10}
      px={5}
    >
      <Stack textAlign={"center"}>
        <Heading fontSize={20}>{title}</Heading>
        <Text fontSize={40} pt={3} fontWeight={"bold"}>
          {number}
        </Text>
      </Stack>
    </Box>
  );
};
StatCard.defaultProps = {
  title: "Verified VASPs",
  number: 248,
};

export default StatCard;
