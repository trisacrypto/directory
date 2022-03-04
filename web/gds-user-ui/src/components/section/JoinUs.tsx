import { ReactNode } from "react";
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Heading,
  Button,
  SimpleGrid,
} from "@chakra-ui/react";
import { Link, useLocation, Navigate, Location } from "react-router-dom";
import { getIcon } from "../Icon";
import { colors } from "../../utils/theme";

const datas = [
  {
    icon: "secure",
    title: "Secure",
    content: (
      <>
        TRISA uses public-key cryptography for encrpyting data in flight and at
        rest.
      </>
    ),
  },
  {
    icon: "network",
    title: "Peer-to-Peer",
    content: (
      <>
        TRISA is a decentralized network where VASPs communicate directly with
        each other.
      </>
    ),
  },
  {
    icon: "opensource",
    title: "Open Source",
    content: <>TRISA is open source and available to implement by any VASPs.</>,
  },
  {
    icon: "plug",
    title: "Interoperable",
    content: (
      <>
        TRISA is designed to be interoperable with other Travel Rule solutuions.
      </>
    ),
  },
];
export default function JoinUsSection() {
  return (
    <Box bg={colors.system.gray} position={"relative"} width="100%">
      <Container
        maxW={"5xl"}
        zIndex={10}
        position={"relative"}
        fontFamily={colors.font}
      >
        <Stack>
          <Stack
            flex={1}
            color={"white"}
            justify={{ lg: "center" }}
            py={{ base: 4, md: 10 }}
          >
            <Box mb={{ base: 10, md: 25 }} color="white">
              <Text fontWeight={600} mb={6} fontSize={"2xl"}>
                Why Join TRISA
              </Text>
              <Text fontSize={"xl"}>
                TRISA is the only global, open source, peer-to-peer and secure
                Travel Rule network. Become a TRISA-certified VASP today. Learn
                how TRISA works.
              </Text>
            </Box>

            <SimpleGrid
              columns={{ base: 1, md: 4 }}
              spacing={8}
              textAlign="center"
            >
              {datas.map((data) => (
                <Box key={data.title} mb={20}>
                  <Text pb={4}>{getIcon(data.icon)}</Text>
                  <Text fontSize={"2xl"} color={"white"} mb={2}>
                    {data.title}
                  </Text>
                  <Text fontSize={"xl"}>{data.content}</Text>
                </Box>
              ))}
            </SimpleGrid>
            <Box alignItems="center" textAlign="center">
              <Button
                bg="#FF7A59"
                w="306px"
                h="64px"
                color="white"
                borderColor="white"
                border="2px"
                _hover={{ bg: "#FF7A77" }}
              >
                Join Today
              </Button>
            </Box>
          </Stack>
          <Flex flex={1} />
        </Stack>
      </Container>
    </Box>
  );
}
