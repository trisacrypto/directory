import React, { FC } from "react";
import {
  Stack,
  Box,
  Text,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
} from "@chakra-ui/react";
import { colors } from "utils/theme";
interface TrisaImplementationReviewProps {
  data: any;
}

const TrisaImplementationReview = (props: TrisaImplementationReviewProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      color={"#252733"}
      maxHeight={367}
      maxWidth={989}
      fontSize={18}
      p={5}
      px={5}
    >
      <Stack>
        <Box display={"flex"} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Section 4: TRISA Implementation</Heading>
          <Button
            bg={colors.system.blue}
            color={"white"}
            height={"34px"}
            _hover={{
              bg: "#10aaed",
            }}
          >
            {" "}
            Edit{" "}
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table
            sx={{
              "td:nth-child(2),td:nth-child(3)": { fontWeight: "bold" },
              Tr: { borderStyle: "hidden" },
            }}
          >
            <Tbody>
              <Tr>
                <Td>TRISA Endpoint</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Certificate Common Name</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
TrisaImplementationReview.defaultProps = {
  data: {},
};
export default TrisaImplementationReview;
