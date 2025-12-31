import React from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Button,
  SimpleGrid,
  Heading,
  Link
} from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { getIcon } from 'components/Icon';
import { colors } from 'utils/theme';
import { Link as RouterLink } from 'react-router-dom';
const datas = [
  {
    icon: 'secure',
    title: t`Secure`,
    content: (
      <Trans id="TRISA uses public-key cryptography for encrypting data in flight and at rest.">
        TRISA uses public-key cryptography for encrypting data in flight and at rest.
      </Trans>
    )
  },
  {
    icon: 'network',
    title: t`Peer-to-Peer`,
    content: (
      <Trans id="TRISA is a decentralized network where VASPs communicate directly with each other.">
        TRISA is a decentralized network where VASPs communicate directly with each other.
      </Trans>
    )
  },
  {
    icon: 'opensource',
    title: t`Open Source`,
    content: (
      <Trans id="TRISA is open source and available to implement by any VASP.">
        TRISA is open source and available to implement by any VASP.
      </Trans>
    )
  },
  {
    icon: 'plug',
    title: t`Interoperable`,
    content: (
      <Trans id="TRISA is designed to be interoperable with other Travel Rule solutions.">
        TRISA is designed to be interoperable with other Travel Rule solutions.
      </Trans>
    )
  }
];
const JoinUsSection: React.FC = () => {
  return (
    <Flex bg={colors.system.gray} position={'relative'} width="100%" py={12}>
      <Container maxW={'5xl'} zIndex={10} position={'relative'} id={'join'}>
        <Stack>
          <Stack flex={1} justify={{ lg: 'center' }}>
            <Box mb={{ base: 10, md: 25 }} color="white">
              <Heading fontWeight={600} pb={6} fontSize={'2xl'} color="#fff">
                <Trans id="Why Join TRISA">Why Join TRISA</Trans>
              </Heading>
              <Text color="#fff" fontSize={{ base: '16px', md: '17px' }}>
                <Trans id="TRISA is a global, open source, peer-to-peer and secure Travel Rule architecture and network designed to be accessible and interoperable. Become a TRISA-certified VASP today.">
                  TRISA is a global, open source, peer-to-peer and secure Travel Rule architecture
                  and network designed to be accessible and interoperable. Become a TRISA-certified
                  VASP today.
                </Trans>{' '}
                <Link
                  isExternal
                  textDecoration={'underline'}
                  href="https://travelrule.io/getting-started-with-trisa/">
                  <Trans id="Learn how TRISA works.">Learn how TRISA works.</Trans>
                </Link>
              </Text>
            </Box>

            <SimpleGrid columns={{ base: 1, md: 4 }} spacing={8} textAlign="center">
              {datas.map((data) => (
                <Box key={data.title} mb={20}>
                  <Text pb={4}>{getIcon(data.icon)}</Text>
                  <Text
                    fontSize={{ base: '16px', md: '17px' }}
                    color={'white'}
                    fontWeight="700"
                    mb={2}>
                    {data.title}
                  </Text>
                  <Text fontSize={{ base: '16px', md: '17px' }} color="#fff">
                    {data.content}
                  </Text>
                </Box>
              ))}
            </SimpleGrid>
            <Stack
              justifyContent={"center"}
              alignItems={"center"}
              direction={['column', 'row']}
              spacing={4}
              >
              <RouterLink to={'/guide'}>
                <Button
                  bg="#FF7A59"
                  color="white"
                  borderColor="white"
                  py={6}
                  width="190px"
                  borderRadius="0px"
                  border="2px solid #fff"
                  _hover={{ bg: '#FF7A77' }}>
                  <Trans id="Join Today">Join Today</Trans>
                </Button>
              </RouterLink>
                <Button
                  bg="#60C4CA"
                  color="white"
                  borderColor="white"
                  py={6}
                  maxWidth={'190px'}
                  width="100%"
                  borderRadius="0px"
                  border="2px solid #fff"
                  _hover={{ bg: '#24a9df' }}
                  as="a"
                  href="https://travelrule.io/members/"
                  target="_blank">
                  <Trans id="View Member List">View Member List</Trans>
                </Button>
            </Stack>
          </Stack>
          <Flex flex={1} />
        </Stack>
      </Container>
    </Flex>
  );
};
export default React.memo(JoinUsSection);
