import React, { FC } from "react";
import {
  HStack,
  Box,
  Icon,
  Text,
  Heading,
  VStack,
  Grid,
} from "@chakra-ui/react";
import { FaCheckCircle, FaDotCircle, FaRegCircle } from "react-icons/fa";
import { GrClose } from "react-icons/gr";
//import { getSteps } from "hooks/useTestnetProgress";
interface ITestnetCertifacteProps {}

const TestNetCertificateProgressBar: FC = (props: ITestnetCertifacteProps) => {
  return (
    <Box
      position={"relative"}
      bg={"white"}
      height={"96px"}
      pt={5}
      boxShadow="0 24px 50px rgba(55,65, 81, 0.25) "
      borderColor={"#C1C9D2"}
      borderRadius={8}
      borderWidth={1}
      mt={10}
      mx={5}
      px={5}
      fontFamily={"Open Sans"}
    >
      <Box pb={2} display={"flex"} justifyContent={"space-between"}>
        <Heading fontSize={20}>Testnet Certification Progress</Heading>
        <Icon as={GrClose} color="#34A853" />
      </Box>

      <Grid templateColumns="repeat(6, 1fr)" gap={2}>
        <Box w="70px" h="1" borderRadius={50} bg="#34A853" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaCheckCircle} color="#34A853" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              Basic Details
            </Text>
          </HStack>
        </Box>

        <Box w="70px" h="1" bg="#5469D4" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaDotCircle} color="#5469D4" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              Legal Person
            </Text>
          </HStack>
        </Box>

        <Box w="70px" h="1" bg="#C1C9D2" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaRegCircle} color="#C1C9D2" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              Contacts
            </Text>
          </HStack>
        </Box>

        <Box w="70px" h="1" bg="#C1C9D2" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaRegCircle} color="#C1C9D2" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              Trisa implementation
            </Text>
          </HStack>
        </Box>

        <Box w="70px" h="1" bg="#C1C9D2" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaRegCircle} color="#C1C9D2" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              TRIXO Questionnaire
            </Text>
          </HStack>
        </Box>

        <Box w="70px" h="1" bg="#C1C9D2" width={"100%"}>
          <HStack>
            <Box pt={3}>
              <Icon as={FaRegCircle} color="#C1C9D2" />
            </Box>
            <Text pt={2} color={"#3C4257"} fontSize={"0.8em"}>
              Submit & Review
            </Text>
          </HStack>
        </Box>
      </Grid>
    </Box>
  );
}; // TestNetCertifateProgressBar

TestNetCertificateProgressBar.defaultProps = {};
export default TestNetCertificateProgressBar;
