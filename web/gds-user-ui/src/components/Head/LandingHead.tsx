import React from 'react';
import { Flex, Heading, Stack, Text, HStack, Button, VStack } from '@chakra-ui/react';

interface LandingHeaderProps {
  hasBtn?: boolean;
  isStartPage?: boolean;
  isHomePage?: boolean;
}
// we should add props to the LandingHead component to allow it to update content dynamically
const LandingHead: React.FC<LandingHeaderProps> = ({ isStartPage, isHomePage, hasBtn }): any => {
  return (
    <Flex
      bgGradient="linear-gradient(90.17deg, rgba(35, 167, 224, 0.85) 33.85%, rgba(27, 206, 159, 0.55) 96.72%);"
      color="white"
      width="100%"
      minHeight={286}
      justifyContent="center"
      direction="column"
      paddingY={{ base: 16, md: 20 }}
      fontSize={'xl'}>
      <Stack textAlign={'center'} color="white" spacing={{ base: 3 }}>
        {isHomePage && (
          <>
            <Heading fontWeight={700} fontSize={{ md: '4xl', sm: '3xl', lg: '5xl' }} color="#fff">
              TRISA Global Directory Service
            </Heading>
            <Text as="h3" fontSize={{ base: '1.2rem', md: 28, sm: 'lg' }} color="#fff">
              Become Travel Rule compliant. <br />
              Apply to Become a TRISA certified Virtual Asset Service Provider.
            </Text>
          </>
        )}
        {isStartPage && (
          <VStack spacing={4}>
            <Heading fontWeight={600} fontSize={{ md: '4xl', sm: '2xl' }} color="#fff">
              Complete TRISA’s VASP Verfication Process
            </Heading>
            <Text fontSize={{ base: 'lg', md: 'xl', sm: 'lg' }} color="#fff" maxW="700px" mx="auto">
              All TRISA members must complete TRISA’s VASP verification and due diligence process to
              become a Verified VASP.
            </Text>
          </VStack>
        )}

        {hasBtn && (
          <>
            <Stack
              justifyContent={'center'}
              alignItems="center"
              pt={4}
              spacing={4}
              direction={['column', 'row']}>
              <Button
                bg={'white'}
                border="3px solid #555151D4"
                color={'black'}
                py={6}
                maxWidth={256}
                width="100%"
                as="a"
                href="/certificate/registration">
                Start Registration
              </Button>
              <Button
                bg={'white'}
                border="3px solid #555151D4"
                color={'black'}
                py={6}
                maxWidth={256}
                width="100%"
                as="a"
                href="/#search">
                Search Directory
              </Button>
            </Stack>
          </>
        )}
      </Stack>
    </Flex>
  );
};

export default LandingHead;
