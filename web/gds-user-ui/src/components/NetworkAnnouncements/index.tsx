import React, { FC } from 'react';
import { Stack, Box, Text, Heading, Flex, VStack } from '@chakra-ui/react';
import AnnouncementCarousel from './caroussels';
import * as Sentry from '@sentry/react';
interface NetworkAnnouncementProps {
  data?: any;
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
      <Sentry.ErrorBoundary
        fallback={
          <Text color={'red'} pt={20}>{`An error has occurred to load annoucements`}</Text>
        }>
        <AnnouncementCarousel announcements={props.data} />
      </Sentry.ErrorBoundary>
    </Flex>
  );
};

export default NetworkAnnouncements;
