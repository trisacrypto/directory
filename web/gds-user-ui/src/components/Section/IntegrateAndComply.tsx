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
  Container,
} from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';
import { colors } from 'utils/theme';
import { Line } from './Line';
import LandingHeader from 'components/Header/LandingHeader';
import Footer from 'components/Footer/LandingFooter';
import { Trans, t } from '@lingui/macro';
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
              <Trans>Integrate with TRISA.</Trans>
            </Heading>
            <Heading
              fontWeight={700}
              fontFamily="Open Sans, sans-serif !important"
              fontSize={{ md: '4xl', sm: '2xl' }}
              color="#fff">
              <Trans>Comply with the Travel Rule</Trans>
            </Heading>
            <Text as="p" mt={2}>
              <Trans>
                Participate in verified VASP-to-VASP Travel Rule exchanges.
              </Trans>
            </Text>
          </VStack>
        </Stack>
      </Flex>
      <Container maxW={'5xl'}>
        <Flex bg={'white'} color={'black'} fontFamily={'Open Sans'}>
          <Stack>
            <Stack flex={1} justify={{ lg: 'center' }} py={6}>
              <Box my={4} color="black">
                <Text fontSize="lg" pb={4}>
                  <Trans>
                  Integrate with <Text as="span" fontWeight={700}>TRISA</Text> using Envoy to begin exchanging Travel Rule compliance data.
                  </Trans>
                </Text>
                <Text fontSize="lg">
                  <Trans>
                  Envoy is an <Text as="span" fontWeight={700}>open source messaging tool</Text> built to help compliance teams handle travel rule exchanges efficiently and securely.
                  It is designed by compliance experts and TRISA engineers to provide simplified data exchanges using the TRISA or TRP protocols.
                  </Trans>
                </Text>
              </Box>
              <Box bg={'gray.100'} p={5}>
                <Text fontSize="lg" color={'black'}>
                  <Trans>
                    VASPs have three options to integrate with TRISA using Envoy.
                  </Trans>
                </Text>
              </Box>
              <Box pt={4}>
                <Flex flexDirection={"column"} gap={"3"} mb={"6"} p={"4"} border="1px solid #00000094" borderRadius={"10px"}>
                  <Heading as="h2" size="md" mb={1}><Trans>1. Open Source/DIY</Trans></Heading>
                  <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Description:</Trans></Text>
                      <Text as="p"><Trans>Envoy is open source (MIT License). Download, install, integrate, host and support your own Envoy node and service.</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Best for:</Trans></Text>
                      <Text as="p"><Trans>Operators that want maximum flexibility and control to customize and have the technical resources to integrate with their systems.</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Est. Time:</Trans></Text>
                      <Text as="p"><Trans>3-6 months</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Cost:</Trans></Text>
                      <Text as="p"><Trans>USD $0</Trans></Text>
                    </Box>
                    <Text as="p" size="sm" fontWeight={600}><Trans>Next Steps:</Trans></Text>
                    <UnorderedList listStylePosition={"inside"}>
                      <ListItem>
                        <Trans>
                        Review <Link isExternal href="https://trisa.dev/envoy/index.html">documentation</Link>
                        </Trans>
                      </ListItem>
                      <ListItem>
                        <Trans>
                        Review <Link isExternal href="https://github.com/trisacrypto/envoy">source code</Link>
                        </Trans>
                      </ListItem>
                      <ListItem>
                        <Trans>
                        Download <Link isExternal href="https://hub.docker.com/r/trisa/trisa">Docker image</Link>
                        </Trans>
                      </ListItem>
                      <ListItem>
                        <Trans>
                          Register with <Link isExternal href="https://trisa.directory/guide">TRISA's GDS</Link>
                        </Trans>
                      </ListItem>
                      <ListItem>
                        <Trans>
                          Start your implementation
                        </Trans>
                      </ListItem>
                    </UnorderedList>
                </Flex>
                <Flex flexDirection={"column"} gap={"3"} mb={"6"} p={"4"} border="1px solid #00000094" borderRadius={"10px"}>
                  <Heading as="h2" size="md"><Trans>2. One-time Setup</Trans></Heading>
                  <Box>
                    <Text as="p" size={"md"} fontWeight={600}><Trans>Description:</Trans></Text>
                    <Text as="p"><Trans>Rotational will install and configure your node for your environment while you host, maintain, and support the node on an ongoing basis.</Trans></Text>
                  </Box>
                  <Box>
                    <Text as="p" size={"md"} fontWeight={600}><Trans>Best for:</Trans></Text>
                    <Text as="p"><Trans>Operators seeking a balance between cost & convenience, with professional implementation ensuring a robust setup without the ongoing costs of managed services.</Trans></Text>
                  </Box>
                  <Box>
                    <Text as="p" size={"md"} fontWeight={600}><Trans>Est. Time:</Trans></Text>
                    <Text as="p"><Trans>2-3 weeks</Trans></Text>
                  </Box>
                  <Box>
                    <Text as="p" size={"md"} fontWeight={600}><Trans>Cost:</Trans></Text>
                    <Text as="p"><Trans>USD USD $2,500 (one-time)</Trans></Text>
                  </Box>
                  <Text as="p" size="sm" fontWeight={600}><Trans>Next Steps:</Trans></Text>
                  <UnorderedList listStylePosition={"inside"}>
                    <ListItem>
                      <Trans>
                        Contact Rotational Labs at <Link href="mailto:support@rotational.io">support@rotational.io</Link> to request information.
                      </Trans>
                    </ListItem>
                    <ListItem>
                      <Trans>
                        Register with <Link isExternal href="https://trisa.directory/guide">TRISA's GDS</Link>
                      </Trans>
                    </ListItem>
                  </UnorderedList>
                </Flex>
                <Flex flexDirection={"column"} gap={"3"} mb={"6"} p={"4"} border="1px solid #00000094" borderRadius={"10px"}>
                    <Heading as="h2" size="md">3. Managed Service</Heading>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Description:</Trans></Text>
                      <Text as="p"><Trans>Rotational will install, configure, host, maintain, and support your node. Includes dedicated, provenance-aware node with regional deployments.</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Best for:</Trans></Text>
                      <Text as="p"><Trans>Compliance teams seeking a turnkey solution with minimal effort while leveraging expert support & maintenance for sustained regulatory adherence.</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Est. Time:</Trans></Text>
                      <Text as="p"><Trans>3-7 days</Trans></Text>
                    </Box>
                    <Box>
                      <Text as="p" size={"md"} fontWeight={600}><Trans>Cost:</Trans></Text>
                      <Text as="p"><Trans>USD $375/month/node (min 6-month commitment)</Trans></Text>
                    </Box>
                    <Text as="p" size="sm" fontWeight={600}><Trans>Next Steps:</Trans></Text>
                      <UnorderedList listStylePosition={"inside"}>
                        <ListItem>
                          <Trans>
                            Contact Rotational Labs at <Link href="mailto:support@rotational.io">support@rotational.io</Link> to request information.
                          </Trans>
                        </ListItem>
                        <ListItem>
                          <Trans>
                            Register with <Link isExternal href="https://trisa.directory/guide">TRISA's GDS</Link>
                          </Trans>
                        </ListItem>
                      </UnorderedList>
                </Flex>
                <Stack direction={['column', 'row']} mt={5} spacing={10}>
                  <Stack bg={'gray.100'} py={5} w="100%">
                    <Line title={t`Need to Learn More?`} fontWeight={'bold'}>
                      <UnorderedList color={'#1F4CED'}>
                        <ListItem>
                          <Link isExternal href="https://trisa.dev/reference/faq/index.html">
                            <Trans>How TRISA Works</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link isExternal href="https://intervasp.org/">
                            <Trans>What is IVMS101?</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link href="https://github.com/trisacrypto/trisa" isExternal>
                            <Trans>TRISA's Github repo</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link href="https://trisa.dev" isExternal>
                            <Trans>Documentation</Trans>
                          </Link>
                        </ListItem>
                        <ListItem>
                          <Link
                            href="https://github.com/trisacrypto/trisa/commit/436fd73fc48973ce09ccbae4260df6213d0c2894"
                            isExternal>
                            <Trans>Reference implementation</Trans>
                          </Link>
                        </ListItem>
                      </UnorderedList>
                    </Line>
                  </Stack>
                </Stack>
              </Box>
              <Stack direction={['column']} pt={5} justifyContent="center">
                <Button as={RouterLink} to={'/auth/register'} alignSelf="center" minWidth={'300px'}>
                  <Trans>Create account</Trans>
                </Button>
                <Text textAlign="center">
                  <Trans>Already have an account?</Trans>{' '}
                  <RouterLink to={'/auth/login'}>
                    <Link color={colors.system.cyan}>
                      {' '}
                      <Trans>Log in.</Trans>
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
