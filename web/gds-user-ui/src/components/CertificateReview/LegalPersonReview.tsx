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
interface LegalSectionProps {
  data: any;
}

const LegalPersonReview = (props: LegalSectionProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      color={"#252733"}
      maxWidth={989}
      fontSize={18}
      p={5}
      px={5}
    >
      <Stack>
        <Box display={"flex"} justifyContent="space-between" pt={4} ml={5}>
          <Heading fontSize={24}>Section 2: Legal Person</Heading>
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
                <Td>Name Identifiers</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td fontStyle={"italic"}>
                  The name and type of name by which the legal person is known.
                </Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Addressess</Td>
                <Td>
                  123 Main Street Legal Person Suite 505 Springfield, CA 90210
                  USA
                </Td>
                <Td>Legal Person</Td>
              </Tr>
              <Tr>
                <Td>Customer Number</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Registration</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>National Identification</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Number</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Identification Type</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Country of Issue</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Reg Authority</Td>
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
LegalPersonReview.defaultProps = {
  data: {},
};
export default LegalPersonReview;
