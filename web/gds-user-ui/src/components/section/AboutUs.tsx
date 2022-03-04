import {
  Container,
  SimpleGrid,
  Image,
  Flex,
  Heading,
  Text,
  Stack,
  StackDivider,
  Fade,
  useColorModeValue
} from '@chakra-ui/react';

import trisaNetworkSvg from '../../assets/trisa_network.svg';
import { colors } from '../../utils/theme';

export default function AboutTrisaSection() {
  return (
    <Container maxW={'5xl'} py={12} fontFamily={colors.font}>
      <SimpleGrid columns={{ base: 1, md: 2 }} spacing={10}>
        <Stack spacing={4}>
          <Heading
            color={useColorModeValue('black', 'white')}
            fontWeight={600}
            fontSize={'2xl'}
            alignSelf={'flex-start'}
            position="relative"
            bottom="0.2rem">
            About Trisa
          </Heading>
          <Text color={useColorModeValue('black', 'white')} fontSize={'lg'}>
            The Travel Rule Information Sharing Alliance (TRISA) is the only global, open source,
            secure, and peer-to-peer protocol for Travel Rule compliance. TRISA helps Virtual Asset
            Service Providers (VASPs) comply with the Travel Rule for cross-boarder cryptocurrency
            transactions. TRISA is designed to be interoperable.
            <br />
            <br />
            TRISA’s Global Directory Service (GDS) is a network of vetted VASPs that can securely
            exchange Travel Rule compliance data. Learn how TRISA works.
          </Text>
          <Stack
            spacing={4}
            divider={
              <StackDivider borderColor={useColorModeValue('gray.100', 'gray.700')} />
            }></Stack>
        </Stack>
        <Flex>
          <Image rounded={'md'} alt={'trisa network'} src={trisaNetworkSvg} />
        </Flex>
      </SimpleGrid>
    </Container>
  );
}
