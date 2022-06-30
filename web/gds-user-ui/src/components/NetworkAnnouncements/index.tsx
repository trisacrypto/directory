import React, { FC } from 'react';
import { Stack, Box, Text, Heading, Flex, HStack } from '@chakra-ui/react';

interface NetworkAnnouncementProps {
  message: string;
  datas?: any;
}
const NetworkAnnouncements = (props: NetworkAnnouncementProps) => {
  return (
    <Flex
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      maxHeight={190}
      fontSize={'18px'}
      bg={'white'}
      p={5}
      mt={10}>
      <Stack>
        <HStack justifyContent="space-between">
          <Heading fontSize={'1.2rem'}> Network announcements</Heading>
          <Box className="arrow">
            <Box className="arrow-top"></Box>
            <Box className="arrow-bottom"></Box>
          </Box>
        </HStack>

        <Text>{props.message}</Text>
      </Stack>
    </Flex>
  );
};

NetworkAnnouncements.defaultProps = {
  message: `Join us on Thursday Jan 28 for the TRISA Working Group call featuring
          guest speaker Jonathon Fishman, Assistant Director, Office of
          Terrorist Financing and Financial Crime at U.S. Department of the
          Treasury.`
};

export default NetworkAnnouncements;
