import React, { FC, useState } from 'react';
import {
  Box,
  IconButton,
  useBreakpointValue,
  Stack,
  Heading,
  Text,
  Container,
  HStack,
  Tag
} from '@chakra-ui/react';
import { BiLeftArrowAlt, BiRightArrowAlt } from 'react-icons/bi';
import Slider from 'react-slick';
import { Trans } from '@lingui/react';

interface AnnouncementProps {
  announcements: any[];
}

const AnnouncementCarousel: FC<AnnouncementProps> = ({ announcements }) => {
  const [currentSlide, setCurrentSlide] = useState(0);
  // slider settings
  const settings = {
    dots: true,
    arrows: false,
    fade: true,
    infinite: true,
    autoplay: false,
    speed: 500,
    autoplaySpeed: 5000,
    slidesToShow: 1,
    slidesToScroll: 1,
    beforeChange: (current: number, next: number) => setCurrentSlide(next),
    afterChange: (current: number) => setCurrentSlide(current)
  };
  // As we have used custom buttons, we need a reference variable to
  // change the state
  const [slider, setSlider] = React.useState<Slider | null>(null);

  // These are the breakpoints which changes the position of the
  // buttons as the screen size changes
  const top = useBreakpointValue({ base: '90%', md: '50%' });
  const side = useBreakpointValue({ base: '30%', md: '40px' });

  //   const datas = [
  //     {
  //       title: 'Upcoming TRISA Working Group Call',
  //       body: 'Join us on Thursday Apr 28 for the TRISA Working Group.',
  //       post_date: '2022-04-20',
  //       author: 'admin@travelrule.io'
  //     },
  //     {
  //       title: 'Routine Maintenance Scheduled',
  //       body: 'The GDS will be undergoing routine maintenance on Apr 7.',
  //       post_date: '2022-04-01',
  //       author: 'admin@travelrule.io'
  //     },
  //     {
  //       title: 'Beware the Ides of March',
  //       body: 'I have a bad feeling about tomorrow.',
  //       post_date: '2022-03-14',
  //       author: 'julius@caesar.com'
  //     }
  //   ];

  return (
    <Box position={'relative'} width={'full'} overflow={'hidden'}>
      {/* CSS files for react-slick */}
      <link
        rel="stylesheet"
        type="text/css"
        // eslint-disable-next-line react/no-unknown-property
        charSet="UTF-8"
        href="https://cdnjs.cloudflare.com/ajax/libs/slick-carousel/1.6.0/slick.min.css"
      />
      <link
        rel="stylesheet"
        type="text/css"
        href="https://cdnjs.cloudflare.com/ajax/libs/slick-carousel/1.6.0/slick-theme.min.css"
      />
      {/* Left Icon */}
      <IconButton
        aria-label="left-arrow"
        variant="ghost"
        position="absolute"
        left={side}
        top={top}
        transform={'translate(0%, -50%)'}
        zIndex={2}
        disabled={currentSlide === 0}
        onClick={() => slider?.slickPrev()}>
        <BiLeftArrowAlt size="40px" />
      </IconButton>
      {/* Right Icon */}
      <IconButton
        aria-label="right-arrow"
        variant="ghost"
        position="absolute"
        right={side}
        top={top}
        transform={'translate(0%, -50%)'}
        zIndex={2}
        disabled={currentSlide === announcements.length - 1}
        onClick={() => slider?.slickNext()}>
        <BiRightArrowAlt size="40px" />
      </IconButton>
      {/* Slider */}
      <Slider {...settings} ref={(sl: any) => setSlider(sl)}>
        {announcements.map((announcement: any, index: any) => (
          <Box key={index} position="relative">
            <Container size="container.lg" position="relative" mt={5}>
              <Stack justifyContent={'center'}>
                <HStack justifyContent={'space-between'}>
                  <Text fontWeight={'bold'} fontSize={'2xl'}>
                    <Trans id="Network Announcement">Network Announcement</Trans>
                  </Text>
                  <Tag bg={'black'} color={'white'}>
                    {announcement.post_date}
                  </Tag>
                  <Tag bg={'blue'} color={'white'}>
                    {announcement.author}
                  </Tag>
                </HStack>
                <Stack>
                  <Heading fontSize={'1rem'}> {announcement.title}</Heading>
                  <Text>{announcement.body}</Text>
                </Stack>
              </Stack>
            </Container>
          </Box>
        ))}
      </Slider>
    </Box>
  );
};

export default AnnouncementCarousel;
