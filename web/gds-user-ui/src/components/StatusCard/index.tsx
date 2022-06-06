import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex, chakra } from '@chakra-ui/react';
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
    <Box
      bg={'white'}
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      // minWidth={250}
      // height={170}
      fontSize={18}
      p={5}
      mt={10}
      px={5}>
      <Stack textAlign={'center'}>
        <chakra.h1 textAlign={'center'} fontSize={20} fontWeight={'bold'}>
          Network Status
        </chakra.h1>
        <Stack
          fontSize={40}
          pt={5}
          alignItems={'center'}
          textAlign={'center'}
          justifyContent={'center'}
          mx={'auto'}>
          {isOnline ? (
            <IoEllipse fontSize="3rem" fill={'#60C4CA'} />
          ) : (
            <IoEllipse fontSize="3rem" fill={'#C4C4C4'} />
          )}
        </Stack>
      </Stack>
    </Box>
  );
};
StatusCard.defaultProps = {
  isOnline: false
};
export default StatusCard;
