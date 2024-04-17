import {
  Stack,
  Box,
  Flex,
  Text,
  Link,
  VStack,
  UnorderedList,
  ListItem,
  Button,
  Heading,
  Container
} from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';
import { colors } from 'utils/theme';
import { Trans } from '@lingui/react';
import { Line } from './Line';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
import { t } from '@lingui/macro';
import LandingBanner from 'components/Banner/LandingBanner';

export default function IntegrateAndComply() {
  return (
    <>
      <LandingHeader />
      <LandingBanner />
      <Flex
        bgGradient="linear-gradient(90.17deg, rgba(35, 167, 224, 0.85) 3.85%, rgba(27, 206, 159, 0.55) 96.72%);"
        color="white"
        width="100%"
        minHeight={286}
        justifyContent="center"
        direction="column"
        paddingY={{ base: 12, md: 16 }}
        fontSize={'xl'}>
        <Stack textAlign={'center'} color="white" spacing={{ base: 3 }}>
          <VStack spacing={1}>
            <Heading
              fontWeight={700}
              fontFamily="Open Sans, sans-serif !important"
              fontSize={{ md: '4xl', sm: '2xl' }}
              color="#fff">
              <Trans id="Integrate with TRISA.">Integrate with TRISA.</Trans>
            </Heading>
            <Heading
              fontWeight={700}
              fontFamily="Open Sans, sans-serif !important"
              fontSize={{ md: '4xl', sm: '2xl' }}
              color="#fff">
              <Trans id="Comply with the Travel Rule">Comply with the Travel Rule</Trans>
            </Heading>
            <Text as="p" mt={2}>
              <Trans id="Participate in verified VASP-to-VASP Travel Rule exchanges.">
                Participate in verified VASP-to-VASP Travel Rule exchanges.
              </Trans>
            </Text>
          </VStack>
        </Stack>
      </Flex>
      <Container maxW={'5xl'}>
        <Flex bg={'white'} color={'black'} fontFamily={'Open Sans'}>
          <Stack>
            <Stack flex={1} justify={{ lg: 'center' }} py={{ base: 4, md: 14 }}>
              <Box pl={5} my={{ base: 4 }} color="black">
                <Text fontSize="lg">
                  <Trans id="Upon verification, integrate with TRISA to begin exchanging Travel Rule compliance data.">
                    Upon verification, integrate with TRISA to begin exchanging Travel Rule
                    compliance
                  </Trans>
                </Text>
              </Box>
              <Box bg={'gray.100'} p={5}>
                <Text fontSize="lg" color={'black'}>
                  <Trans id="VASPs have two options to integrate with TRISA.">
                    VASPs have two options to integrate with TRISA.
                  </Trans>
                </Text>
              </Box>
              <Box mt={20} pt={10}>
                <Stack
                  spacing={{ base: 20, md: 0 }}
                  display={{ md: 'grid' }}
                  gridTemplateColumns={{ md: 'repeat(2,1fr)' }}
                  color={'black'}
                  gridColumnGap={{ md: 20, lg: 54 }}
                  gridRowGap={{ md: 6 }}
                  position="relative">
                  <Line title={t`Option 1. Set Up Your Own TRISA Node`}>
                    <Trans
                      id="Since TRISA is an open source, peer-to-peer Travel Rule solution, VASPs can set
                    up and maintain their own TRISA server to exhange encrypted Travel Rule
                    compliance data. TRISA maintains an">
                      Since TRISA is an open source, peer-to-peer Travel Rule solution, VASPs can
                      set up and maintain their own TRISA server to exhange encrypted Travel Rule
                      compliance data. TRISA maintains an
                    </Trans>{' '}
                    <Link href="https://github.com/trisacrypto/trisa" isExternal color="link">
                      <Trans id="GitHub repository">GitHub repository</Trans>
                    </Link>{' '}
                    <Trans
                      id="with detailed documentation, a reference implemenation, and “robot” VASPs for
                    testing purposes.">
                      with detailed documentation, a reference implemenation, and “robot” VASPs for
                      testing purposes.
                    </Trans>
                  </Line>
                  <Line title={t`Option 2. Use a 3rd-party Solution`}>
                    <Trans
                      id="TRISA is designed to be interoperable. There are several Travel Rule solutions
                    providers available on the market. If you are a customer, work with them to
                    integrate TRISA into your Travel Rule compliance workflow.">
                      TRISA is designed to be interoperable. There are several Travel Rule solutions
                      providers available on the market. If you are a customer, work with them to
                      integrate TRISA into your Travel Rule compliance workflow.
                    </Trans>
                  </Line>
                  <Line title={t`How to set up your own node?`} fontStyle="italic">
                    <Trans
                      id="Talk to a member of your technical team to determine the requirements and
                    resources to integrate TRISA with your system. Have members of your technical
                    team integrate your systems with TRISA. Or work with a solutions provider that
                    can help your VASP set up your TRISA server and maintain it.">
                      Talk to a member of your technical team to determine the requirements and
                      resources to integrate TRISA with your system. Have members of your technical
                      team integrate your systems with TRISA. Or work with a solutions provider that
                      can help your VASP set up your TRISA server and maintain it.
                    </Trans>
                  </Line>
                  <Line title={t`3rd Party Travel Rule Providers`}>
                    <UnorderedList color={'#1F4CED'}>
                      <ListItem>
                        <Link href="https://ciphertrace.com/travel-rule-compliance/" isExternal>
                          <Trans id="CipherTrace Traveler">CipherTrace Traveler</Trans>
                        </Link>
                      </ListItem>
                      <ListItem>
                        <Link href="https://sygna.io" isExternal>
                          <Trans id="Synga Bridge">Synga Bridge</Trans>
                        </Link>
                      </ListItem>
                    </UnorderedList>
                  </Line>
                </Stack>
                <Stack direction={['column', 'row']} mt={5} spacing={10}>
                  <Stack bg={'gray.100'} py={5} w="100%">
                    <Line title={t`Open Source Resources`} fontWeight={'bold'}>
                      <UnorderedList color={'#1F4CED'}>
                        <ListItem>
                          <Link href="https://github.com/trisacrypto/trisa" isExternal>
                            <Trans id="TRISA’s Github repo">TRISA’s Github repo</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link href="https://trisa.dev/" isExternal>
                            <Trans id="Documentation">Documentation</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link
                            href="https://github.com/trisacrypto/trisa/commit/436fd73fc48973ce09ccbae4260df6213d0c2894"
                            isExternal>
                            <Trans id="Reference implementation">Reference implementation</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link href="https://vaspbot.net/" isExternal>
                            <Trans id="Meet Alice VASP, Bob VASP, and “Evil” VASP">
                              Meet Alice VASP, Bob VASP, and “Evil” VASP
                            </Trans>
                          </Link>
                        </ListItem>
                      </UnorderedList>
                    </Line>
                  </Stack>
                  <Stack bg={'gray.100'} py={5} w="100%">
                    <Line title={t`Need to Learn More?`} fontWeight={'bold'}>
                      <UnorderedList color={'#1F4CED'}>
                        <ListItem>
                          <Link isExternal href="https://trisa.io/getting-started-with-trisa/">
                            <Trans id="Learn How TRISA Works">Learn How TRISA Works</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link isExternal href="https://intervasp.org/">
                            <Trans id="What is IVMS101?">What is IVMS101?</Trans>
                          </Link>
                        </ListItem>
                      </UnorderedList>
                    </Line>
                  </Stack>
                </Stack>
              </Box>
              <Stack direction={['column']} pt={5} justifyContent="center">
                <Button as={RouterLink} to={'/auth/register'} alignSelf="center" minWidth={'300px'}>
                  <Trans id="Create account">Create account</Trans>
                </Button>
                <Text textAlign="center">
                  <Trans id="Already have an account?">Already have an account?</Trans>{' '}
                  <RouterLink to={'/auth/login'}>
                    <Link color={colors.system.cyan}>
                      {' '}
                      <Trans id="Log in.">Log in.</Trans>
                    </Link>
                  </RouterLink>
                </Text>
              </Stack>
            </Stack>
          </Stack>
        </Flex>
      </Container>
      <Footer />
    </>
  );
}
