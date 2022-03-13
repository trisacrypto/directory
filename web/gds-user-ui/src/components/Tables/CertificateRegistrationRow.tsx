import React from "react";
import {
  Tr,
  Td,
  Flex,
  Text,
  Icon,
  IconButton,
  Box,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  useColorModeValue,
} from "@chakra-ui/react";

import { BsThreeDots } from "react-icons/bs";

interface CertificateRegistrationRowProps {
  name: string;
  section: string;
  description: string;
  status: string | null;
}

const getBackgroundByStatusCode = (status: string | null) => {
  if (!status) {
    return "#555151";
  }
};

const CertificateRegistrationRow = (props: CertificateRegistrationRowProps) => {
  const { section, name, description, status } = props;
  const textColor = useColorModeValue("#858585", "white");
  return (
    <Tr
      border="1px solid #23A7E0"
      borderRadius={100}
      sx={{
        td: {
          height: "66px",
          borderTop: "1px solid #23A7E0",
          borderBottom: "1px solid #23A7E0",
        },
        "td:first-child": {
          border: "1px solid #23A7E0",
          borderLeftRadius: 100,
          borderRight: "none",
          textAlign: "center",
        },
        "td:nth-child(3)": {},
        "td:last-child": {
          textAlign: "center",
          border: "1px solid #23A7E0",
          borderRightRadius: 100,
          borderLeft: "none",
        },
      }}
    >
      <Td>
        <Text fontSize="md" color={textColor} pb=".5rem">
          {section}
        </Text>
      </Td>
      <Td minWidth={{ sm: "250px" }} pl="0px">
        <Flex alignItems="center" py=".8rem" minWidth="100%" flexWrap="nowrap">
          <Text
            fontSize="md"
            color={"#1F4CED"}
            fontWeight="bold"
            minWidth="100%"
          >
            {name}
          </Text>
        </Flex>
      </Td>
      <Td>
        <Text fontSize="md" color={textColor} pb=".5rem">
          {description}
        </Text>
      </Td>
      <Td>
        <Text fontSize="md" color={textColor} fontWeight="bold" pb=".5rem">
          <Box
            textAlign="center"
            fontSize={"sm"}
            width="96px"
            alignItems="center"
            height="27px"
            borderRadius="215"
            color="white"
            py={1}
            background={getBackgroundByStatusCode(status)}
          >
            <Text fontSize={"sm"}> {status || "Incomplete"}</Text>
          </Box>
        </Text>
      </Td>

      <Td>
        <Menu>
          <MenuButton
            as={IconButton}
            icon={<BsThreeDots />}
            borderRadius={50}
            bg={"transparent"}
          />
          <MenuList>
            <MenuItem>Edit</MenuItem>
          </MenuList>
        </Menu>
      </Td>
    </Tr>
  );
};

export default CertificateRegistrationRow;
