import React from 'react';
import { Flex, Heading, Stack, Text, Button, VStack, Box } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

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
      minHeight={286}
      justifyContent="center"
      direction="column"
      paddingY={{ base: 16, md: 20 }}
      fontSize={'xl'}
      style={{ marginTop: 0 }}
      >
      <Stack textAlign={'center'} spacing={{ base: 3 }}>
        {isHomePage && (
          <Box width={'80%'} mx="auto">
            <Heading color="#fff" fontWeight={700} fontSize={{ md: '4xl', sm: '3xl', lg: '5xl' }}>
              <Trans id="TRISA Global Directory Service">TRISA Global Directory Service</Trans>
            </Heading>
            <Text
              as="h3"
              fontFamily={'Open Sans, sans-serif !important'}
              fontSize={{ base: '1rem', md: 24, sm: 'lg' }}
              color="#fefefe">
              <Trans id="Become Travel Rule compliant.">Become Travel Rule compliant.</Trans> <br />
              <Trans id="Apply to Become a TRISA-certified Virtual Asset Service Provider.">
                Apply to Become a TRISA-certified Virtual Asset Service Provider.
              </Trans>
            </Text>
          </Box>
        )}
        {isStartPage && (
          <VStack spacing={4}>
            <Heading
              fontWeight={700}
              fontFamily="Open Sans, sans-serif !important"
              fontSize={{ md: '4xl', sm: '2xl' }}
              color="#fff">
              <Trans id="Complete TRISA’s VASP Verfication Process">
                Complete TRISA’s VASP Verification Process
              </Trans>
            </Heading>
            <Text
              fontSize={{ base: 'lg', md: 'xl', sm: 'lg' }}
              fontWeight="400"
              color="#fff"
              maxW="700px"
              mx="auto">
              <Trans id="All TRISA members must complete TRISA’s VASP verification and due diligence process to become a Verified VASP.">
                All TRISA members must complete TRISA’s VASP verification and due diligence process
                to become a Verified VASP.
              </Trans>
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
                as="a"
                _hover={{
                  background: '#24a9df',
                  transition: '0.15s all ease'
                }}
                href="/guide">
                <Trans id="Start Registration">Start Registration</Trans>
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
                as="a"
                href="/#search">
                <Trans id="Directory Lookup">Directory Lookup</Trans>
              </Button>
            </Stack>
          </>
        )}
      </Stack>
    </Flex>
  );
};

export default LandingHead;
