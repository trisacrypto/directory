import React from "react";
import {
  Stack,
  Box,
  Flex,
  Text,
  Link,
  chakra,
  FlexProps,
  StyleProps,
  useColorModeValue,
  UnorderedList,
  ListItem,
  Button,
} from "@chakra-ui/react";

import { colors } from "utils/theme";

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
  };

interface ILineProps {
  children?: React.ReactNode;
  title?: string;
}

const Line: React.FC<Props> = ({ children, title, ...rest }: any) => {
  return (
    <Flex>
      <Flex shrink={0}>
        <Flex
          rounded="md"
          bg={useColorModeValue("brand.500", "")}
          color="white"
        ></Flex>
      </Flex>
      <Box ml={4}>
        <chakra.dt fontSize="lg" fontWeight="medium" lineHeight="6" {...rest}>
          {title}
        </chakra.dt>
        <chakra.dd mt={2}>{children}</chakra.dd>
      </Box>
    </Flex>
  );
};
export default function VaspVerification() {
  return (
    <Flex bg={"white"} color={"black"} fontFamily={"Open Sans"}>
      <Stack>
        <Stack
          flex={1}
          justify={{ lg: "center" }}
          py={{ base: 4, md: 20 }}
          px={{ base: 10, md: 55 }}
        >
          <Box my={{ base: 4 }} color="black">
            <Text fontFamily={"heading"} fontWeight={700} fontSize={"xl"}>
              Upon verification, VASPs must integrate with TRISA to begin
              exchanging Travel Rule data with other verified TRISA members.
            </Text>
          </Box>
          <Box bg={"gray.100"} p={5}>
            <Text fontSize={"xl"} color={"black"}>
              VASPs have two options to integrate with TRISA.
            </Text>
          </Box>
          <Box mt={20} pt={10}>
            <Stack
              spacing={{ base: 20, md: 0 }}
              display={{ md: "grid" }}
              gridTemplateColumns={{ md: "repeat(2,1fr)" }}
              color={"black"}
              gridColumnGap={{ md: 20, lg: 80 }}
              gridRowGap={{ md: 10 }}
            >
              <Line
                title="Option 1. Set Up Your Own TRISA Node"
                fontWeight={"bold"}
              >
                Since TRISA is an open source, peer-to-peer Travel Rule
                solution, VASPs can set up and maintain their own TRISA server
                to exhange encrypted Travel Rule compliance data. TRISA
                maintains an GitHub repository with detailed documentation, a
                reference implemenation, and “robot” VASPs for testing purposes.
              </Line>

              <Line
                title="Option 2. Use a 3rd-party Solution"
                fontWeight={"bold"}
              >
                TRISA is designed to be interoperable. There are several Travel
                Rule solutions providers available on the market. If you are a
                customer, work with them to integrate TRISA into your Travel
                Rule compliance workflow.
              </Line>

              <Line title="How to set up your own node?" fontWeight={"bold"}>
                Talk to a member of your technical team to determine the
                requirements and resources to integrate TRISA with your system.
                Have members of your technical team integrate your systems with
                TRISA. Or work with a solutions provider that can help your VASP
                set up your TRISA server and maintain it.
              </Line>
              <Line title="3rd Party Travel Rule Providers" fontWeight={"bold"}>
                <UnorderedList>
                  <ListItem color={"#1F4CED"}>
                    <Link>CipherTrace</Link>
                  </ListItem>
                  <ListItem color={"#1F4CED"}>
                    <Link>Synga Bridge</Link>
                  </ListItem>
                  <ListItem>
                    <Link color={"#1F4CED"}>NotaBene</Link> (not interoperable)
                  </ListItem>
                  <ListItem>
                    <Link color={"#1F4CED"}>OpenVASP</Link> (not interoperable)
                  </ListItem>
                </UnorderedList>
              </Line>

              <Stack mt={20} bg={"gray.100"} py={5}>
                <Line title="Open Source Resources" fontWeight={"bold"}>
                  <UnorderedList color={"#1F4CED"}>
                    <ListItem>
                      <Link>TRISA’s Github repo</Link>
                    </ListItem>
                    <ListItem>
                      <Link>Documentation</Link>
                    </ListItem>
                    <ListItem>
                      <Link>Reference implementation</Link>
                    </ListItem>
                    <ListItem>
                      <Link>Meet Alice VASP, Bob VASP, and “Evil” VASP</Link>
                    </ListItem>
                  </UnorderedList>
                </Line>
              </Stack>
              <Stack mt={20} bg={"gray.100"} py={5}>
                <Line title="Need to Learn More?" fontWeight={"bold"}>
                  <UnorderedList color={"#1F4CED"}>
                    <ListItem>
                      <Link>How TRISA Works</Link>
                    </ListItem>
                    <ListItem>
                      <Link>What is IVMS101?</Link>
                    </ListItem>
                  </UnorderedList>
                </Line>
              </Stack>
            </Stack>
          </Box>
          <Stack direction={["column", "row"]} pt={10} mx={10}>
            <Box>
              <Button
                bg={colors.system.blue}
                color={"white"}
                _hover={{
                  bg: "#10aaed",
                }}
                _focus={{
                  borderColor: "transparent",
                }}
              >
                Back to Getting Started
              </Button>
            </Box>
          </Stack>
        </Stack>
      </Stack>
    </Flex>
  );
}
