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
      bgGradient="linear-gradient(90.17deg, rgba(35, 167, 224, 0.85) 3.85%, rgba(27, 206, 159, 0.55) 96.72%);"
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
            <Heading
              fontWeight={700}
              fontSize={{ md: '4xl', sm: '3xl', lg: '5xl' }}
              color="#fefefe">
              TRISA Global Directory Service
            </Heading>
            <Text
              as="h3"
              fontFamily={'Open Sans, sans-serif !important'}
              fontSize={{ base: '1rem', md: 24, sm: 'lg' }}
              color="#fefefe">
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
                bg={'transparent'}
                borderRadius="0px"
                border="2px solid #fff"
                color={'#fff'}
                py={6}
                maxWidth={'180px'}
                width="100%"
                as="a"
                _hover={{
                  background: '#24a9df',
                  transition: '0.15s all ease'
                }}
                href="/certificate/registration">
                Start Registration
              </Button>
              <Button
                bg={'transparent'}
                border="2px solid #fff"
                borderRadius="0px"
                color={'#fff'}
                _hover={{
                  background: '#24a9df',
                  transition: '0.15s all ease'
                }}
                py={6}
                maxWidth={'180px'}
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
