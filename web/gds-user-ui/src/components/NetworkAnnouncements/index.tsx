import React, { FC } from 'react';
import { Stack, Box, Text, Heading, Flex, VStack } from '@chakra-ui/react';
import AnnouncementCarousel from './caroussels';
import * as Sentry from '@sentry/react';
interface NetworkAnnouncementProps {
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
      <Sentry.ErrorBoundary
        fallback={
          <Text color={'red'} pt={20}>{`An error has occurred to load annoucements`}</Text>
        }>
        <AnnouncementCarousel announcements={props.datas} />
      </Sentry.ErrorBoundary>
    </Flex>
  );
};

NetworkAnnouncements.defaultProps = {
  datas: [
    {
      title: 'Upcoming TRISA Working Group Call',
      body: 'Join us on Thursday Apr 28 for the TRISA Working Group.',
      post_date: '2022-04-20',
      author: 'admin@trisa.io'
    },
    {
      title: 'Routine Maintenance Scheduled',
      body: 'The GDS will be undergoing routine maintenance on Apr 7.',
      post_date: '2022-04-01',
      author: 'admin@trisa.io'
    },
    {
      title: 'Beware the Ides of March',
      body: 'I have a bad feeling about tomorrow.',
      post_date: '2022-03-14',
      author: 'julius@caesar.com'
    }
  ]
};

export default NetworkAnnouncements;
