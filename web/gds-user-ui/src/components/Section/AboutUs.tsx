import {
  Container,
  SimpleGrid,
  Image,
  Flex,
  Text,
  Stack,
  StackDivider,
  useColorModeValue,
  Link
} from '@chakra-ui/react';

import trisaNetworkSvg from 'assets/trisa_network.svg';
import { colors } from 'utils/theme';
import { Trans } from '@lingui/react';

export default function AboutTrisaSection() {
  return (
    <Flex>
      <Container maxW={'5xl'} py={12} fontFamily={colors.font} id={'about'}>
        <SimpleGrid columns={{ base: 1, md: 2 }} spacing={10}>
          <Stack spacing={4}>
            <Text
              color={useColorModeValue('black', 'white')}
              fontWeight={600}
              pb={6}
              fontSize={'2xl'}
              fontFamily="Roboto Slab">
              <Trans id={'About TRISA'}>About TRISA</Trans>
            </Text>

            <Text
              color={useColorModeValue('black', 'white')}
              fontSize={{ base: '16px', md: '17px' }}>
              <Link isExternal color={colors.system.link} href={'https://trisa.io'}>
                The Travel Rule Information Sharing Architecture (TRISA){' '}
              </Link>
              is a global, open source, secure, and peer-to-peer protocol for{' '}
              <Link
                isExternal
                color={colors.system.link}
                href="https://www.fatf-gafi.org/publications/?hf=10&b=0&s=desc(fatf_releasedate)">
                Travel Rule
              </Link>{' '}
              compliance. TRISA helps Virtual Asset Service Providers (VASPs) comply with the Travel
              Rule for cross-border cryptocurrency transactions. TRISA is designed to be
              interoperable.
              <br />
              <br />
              TRISA’s Global Directory Service (GDS) is a network of vetted VASPs that can securely
              exchange Travel Rule compliance data with each other.
            </Text>
            <Stack
              spacing={4}
              divider={
                <StackDivider borderColor={useColorModeValue('gray.100', 'gray.700')} />
              }></Stack>
          </Stack>
          <Flex>
            <Image alt={'trisa network'} src={trisaNetworkSvg} />
          </Flex>
        </SimpleGrid>
      </Container>
    </Flex>
  );
}
