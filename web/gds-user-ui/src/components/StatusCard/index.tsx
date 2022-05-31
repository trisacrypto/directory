import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex } from '@chakra-ui/react';
import { IoEllipse } from 'react-icons/io5';
interface StatusCardProps {
  isOnline: boolean;
}

const StatusCard = ({ isOnline }: StatusCardProps) => {
  //   return (
  //     <Box
  //       border="1px solid #DFE0EB"
  //       fontFamily={"Open Sans"}
  //       color={"#252733"}
  //       height={167}
  //       maxWidth={451}
  //       fontSize={18}
  //       p={5}
  //       mt={10}
  //       px={5}
  //     >
  //       <Stack>
  //         <Heading fontSize={20}>Certification Status</Heading>
  //         <HStack spacing={10}>
  //           <Text>Testnet</Text>
  //           <Text>{testnetstatus}</Text>
  //         </HStack>
  //         <HStack spacing={8}>
  //           <Text>Mainnet</Text>
  //           <Text>{mainnetstatus}</Text>
  //         </HStack>
  //       </Stack>
  //     </Box>
  //   );
  // };
  // StatusCard.defaultProps = {
  //   testnetstatus: "In progress",
  //   mainnetstatus: "Not Eligible yet ",
  // };
  return (
    <Flex
      bg={'white'}
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      height={167}
      maxWidth={246}
      fontSize={18}
      p={5}
      px={5}>
      <Box textAlign={'center'}>
        <Heading fontSize={20}>Network Status</Heading>
        <Stack
          fontSize={40}
          pt={3}
          fontWeight={'bold'}
          textAlign={'center'}
          justifyContent={'center'}>
          {isOnline ? (
            <IoEllipse fontSize="2rem" fill={'#34A853'} />
          ) : (
            <IoEllipse fontSize="2rem" fill={'#C4C4C4'} />
          )}
        </Stack>
      </Box>
    </Flex>
  );
};
StatusCard.defaultProps = {
  isOnline: false
};
export default StatusCard;
