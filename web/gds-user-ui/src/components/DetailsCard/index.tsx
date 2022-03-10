import React from "react";
import {
  Stack,
  Box,
  Text,
  Heading,
  UnorderedList,
  ListItem,
  HStack,
} from "@chakra-ui/react";

enum DetailType {
  ORG = "org",
  CERT = "cert",
}
interface DetailsCardProps {
  type: string;
  data: any;
}

const DetailsCard = ({ type, title, data }: DetailsCardProps) => {
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={"Open Sans"}
      color={"#252733"}
      height={248}
      maxWidth={473}
      fontSize={18}
      p={5}
      mt={10}
      px={5}
    >
      <Stack>
        {type === DetailType.ORG ? (
          <Stack>
            <Heading fontSize={20}>Organizational Details</Heading>
            <UnorderedList p={5} mt={10} px={5}>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>TRISA Member ID:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>TRISA Verification:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>Country:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
            </UnorderedList>
          </Stack>
        ) : (
          <Stack>
            <Heading fontSize={20}>Certificate Details</Heading>
            <UnorderedList p={5} mt={10} px={5}>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>Organization:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>Issue Date:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>Expiry Date:</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
              <ListItem>
                <HStack justifyContent={"space-between"}>
                  <Text>TRISA Identity Signature::</Text>
                  <Text></Text>
                </HStack>
              </ListItem>
            </UnorderedList>
          </Stack>
        )}
      </Stack>
    </Box>
  );
};
DetailsCard.defaultProps = {
  type: "org",
  data: [],
};

export default DetailsCard;
