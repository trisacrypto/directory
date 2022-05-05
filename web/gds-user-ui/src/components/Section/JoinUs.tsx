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
import { getIcon } from 'components/Icon';
import { colors } from 'utils/theme';

const datas = [
  {
    icon: 'secure',
    title: 'Secure',
    content: <>TRISA uses public-key cryptography for encrpyting data in flight and at rest.</>
  },
  {
    icon: 'network',
    title: 'Peer-to-Peer',
    content: <>TRISA is a decentralized network where VASPs communicate directly with each other.</>
  },
  {
    icon: 'opensource',
    title: 'Open Source',
    content: <>TRISA is open source and available to implement by any VASP.</>
  },
  {
    icon: 'plug',
    title: 'Interoperable',
    content: <>TRISA is designed to be interoperable with other Travel Rule solutuions.</>
  }
];
export default function JoinUsSection() {
  return (
    <Flex bg={colors.system.gray} position={'relative'} width="100%" py={12}>
      <Container maxW={'5xl'} zIndex={10} position={'relative'} id={'join'}>
        <Stack>
          <Stack flex={1} justify={{ lg: 'center' }}>
            <Box mb={{ base: 10, md: 25 }} color="white">
              <Heading fontWeight={600} pb={6} fontSize={'2xl'} color="#fff">
                Why Join TRISA
              </Heading>
              <Text color="#fff" fontSize={{ base: '16px', md: '17px' }}>
                TRISA is a global, open source, peer-to-peer and secure Travel Rule architecture and
                network designed to be accessible and interoperable. Become a TRISA-certified VASP
                today.{' '}
                <Link isExternal href="htps://trisa.io/getting-started-with-trisa/">
                  Learn how TRISA works.
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
            <Box alignItems="center" textAlign="center">
              <Button
                bg="#FF7A59"
                color="white"
                borderColor="white"
                py={6}
                maxWidth={'190px'}
                width="100%"
                borderRadius="0px"
                border="2px solid #fff"
                _hover={{ bg: '#FF7A77' }}
                as={'a'}
                href={'/getting-started'}>
                Join Today
              </Button>
            </Box>
          </Stack>
          <Flex flex={1} />
        </Stack>
      </Container>
    </Flex>
  );
}
