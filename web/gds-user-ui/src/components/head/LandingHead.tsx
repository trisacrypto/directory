import React from 'react';
import { Flex, Heading, Stack, Text, HStack, Button } from '@chakra-ui/react';

interface LandingHeaderProps {
  hasBtn?: boolean;
  isStartPage?: boolean;
  isHomePage?: boolean;
}
// we should add props to the LandingHead component to allow it to update content dynamically
const LandingHead: React.FC<LandingHeaderProps> = ({ isStartPage, isHomePage, hasBtn }): any => {
  return (
    <Flex
      bgGradient="linear(270deg,#24a9df,#1aebb4)"
      color="white"
      width="100%"
      height={286}
      justifyContent="center"
      direction="column"
      padding={4}
      fontSize={'xl'}>
      <Stack textAlign={'center'} color="white" spacing={{ base: 3 }}>
        {isHomePage && (
          <>
            <Heading fontWeight={600} fontSize={{ md: '4xl', sm: 'xl', lg: '3xl' }}>
              TRISA Global Directory Service
            </Heading>
            <Text fontSize={{ base: '30px', md: '2xl', sm: 'lg' }}>
              Become Travel Rule compliant. <br />
              Apply to Become a TRISA certified Virtual Asset Service Provider.
            </Text>
          </>
        )}
        {isStartPage && (
          <>
            <Heading fontWeight={600} fontSize={{ md: '4xl', sm: '2xl' }}>
              Complete TRISA’s VASP Verfication Process
            </Heading>
            <Text fontSize={{ base: '30px', md: '2xl', sm: 'lg' }}>
              All TRISA members must complete TRISA’s VASP verification <br />
              and due diligence process to become a Verified VASP.
            </Text>
          </>
        )}

        {hasBtn && (
          <>
            <Stack justifyContent={'center'} pt={4} spacing={4} direction={['column', 'row']}>
              <Button
                bg={'white'}
                color={'black'}
                p={2}
                minWidth={306}
                minHeight={65}
                as="a"
                href="/#join">
                {' '}
                Why Join TRISA{' '}
              </Button>
              <Button
                bg={'white'}
                color={'black'}
                p={2}
                minWidth={306}
                minHeight={65}
                as="a"
                href="/#search">
                {' '}
                Search Directory{' '}
              </Button>
            </Stack>
          </>
        )}
      </Stack>
    </Flex>
  );
};

export default LandingHead;
