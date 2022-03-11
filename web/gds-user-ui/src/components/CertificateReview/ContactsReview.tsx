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
interface ContactsProps {
  data: any;
}

const ContactsReview = (props: ContactsProps) => {
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
          <Heading fontSize={24}>Section 3: Contacts</Heading>
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
                <Td>Technical Contact</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Compliance/ Legal Contact</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Administrative Contact</Td>
                <Td></Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td>Billing Contact</Td>
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
ContactsReview.defaultProps = {
  data: {},
};
export default ContactsReview;
