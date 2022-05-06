import React from 'react';
import {
  Stack,
  Container,
  Box,
  Flex,
  Text,
  Link,
  chakra,
  FlexProps,
  StyleProps,
  useColorModeValue,
  UnorderedList,
  ListItem,
  Button,
  GridItem,
  useMediaQuery
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
    colSpan: number;
    dataContent?: string;
  };

interface ILineProps {
  children?: React.ReactNode;
  title?: string;
}

const Line: React.FC<Props> = ({ children, colSpan, title, dataContent, ...rest }: any) => {
  const [isMobile] = useMediaQuery('(max-width: 768px)');
  return (
    <GridItem ml={4} colSpan={colSpan}>
      <chakra.dt fontSize="lg" fontWeight="medium" lineHeight="6" mt={2} {...rest}>
        {title}
      </chakra.dt>
      <chakra.dd
        mt={2}
        data-content={dataContent}
        mb={{ base: 4, sm: 0 }}
        fontSize={{ base: '1rem' }}
        padding={2}
        {...(isMobile && {
          _before: {
            content: 'attr(data-content)',
            display: 'block',
            fontWeight: 'bold',
            fontSize: 'lg'
          }
        })}>
        {children}
      </chakra.dd>
    </GridItem>
  );
};
export default function VaspVerification() {
  return (
    <Flex color={'black'} fontFamily={'Open Sans'} fontSize={'xl'}>
      <Container maxW={'5xl'}>
        <Stack flex={1} justify={{ lg: 'center' }}>
          <Box my={{ base: 4 }} color="black">
            <Text fontFamily={'heading'} fontWeight={700} fontSize={'xl'}>
              TRISA members must complete a comprehensive multi-part verification form and due
              diligence process. Once verified, TRISA will issue TestNet and MainNet certificates
              for secure Travel Rule compliance.
            </Text>
          </Box>
          <Box bg={'gray.100'} p={5}>
            <Text fontSize={'1rem'} color={'black'}>
              TRISA’s verification form includes five sections and may require information from
              several parties in your organization.
            </Text>
          </Box>
          <Box mx={20} pt={10}>
            <Box
              display={{ md: 'grid' }}
              gridTemplateColumns={{ md: 'repeat(5,1fr)' }}
              color={'black'}
              gap={10}
              alignItems="center">
              <>
                <Line title="Sections & Details" fontWeight={'bold'} colSpan={3}>
                  {''}
                </Line>
                <Line title="Who to Ask" fontWeight={'bold'} colSpan={2}>
                  {''}
                </Line>
              </>

              <Line title="1 Basic Details" fontWeight={'bold'} colSpan={3}>
                Information about the VASP such as website, incorporation date, business and VASP
                category.
              </Line>

              <Line colSpan={2} dataContent="Who to ask">
                Business or Compliance Office
              </Line>

              <Line title="2 Legal Person" fontWeight={'bold'} colSpan={3}>
                Information that identifies your organization as a Legal Person. This section
                represents the IVMS 101 data structure for legal persons and is strongly suggested
                for use as KYC information exchanged in TRISA transfers.
              </Line>
              <Line colSpan={2} dataContent="Who to ask">
                Business or Compliance Office
              </Line>

              <Line title="3 Contacts" fontWeight={'bold'} colSpan={3}>
                Contact information for representatives of your organization. Contacts include
                Technical, Legal/Compliance, Administrative, and Billing persons.
              </Line>
              <Line colSpan={2} dataContent="Who to ask">
                Business or Compliance Office
              </Line>

              <Line title="4 TRISA Implementation" fontWeight={'bold'} colSpan={3}>
                Technical information about your endpoint for certificate issuance. Each VASP is
                required to establish a TRISA endpoint for inter-VASP communication.
              </Line>
              <Line colSpan={2} dataContent="Who to ask">
                Technical Officer
              </Line>
              <Line title="5 TRIXO Questionnaire" fontWeight={'bold'} colSpan={3}>
                Information to ensure that required compliance information exchanges are conducted
                correctly and safely. This includes information about jurisdiction and national
                regulator, Customer Due Diligence(CDD) and Travel Rule policies, and data protection
                policies.
              </Line>
              <Line colSpan={2} dataContent="Who to ask">
                Compliance Officer
              </Line>
            </Box>
            <Box
              display={{ md: 'grid' }}
              gridTemplateColumns={{ md: 'repeat(5,1fr)' }}
              color={'black'}
              gridColumnGap={10}
              gridRowGap={10}>
              <GridItem ml={4} colSpan={3} bg={'#eee'} mt={5} p={2}>
                <chakra.dt fontSize="lg" fontWeight="bold" lineHeight="6">
                  Final Confirmation
                </chakra.dt>
                <chakra.dd mt={2} fontSize="1rem">
                  For MainNet certificate requests, a member of TRISA’s verification team will
                  review your submission and conduct a final due diligence phone call for physical
                  verification. When physical verification is complete, TRISA will issue MainNet
                  certificates. Requests for TestNet certificates do not require physical
                  verification.
                </chakra.dd>
              </GridItem>
              <GridItem ml={4} colSpan={2} bg={'#eee'} mt={5} p={2}>
                <chakra.dt fontSize="lg" fontWeight="bold" lineHeight="6">
                  Need to Learn More?
                </chakra.dt>
                <chakra.dd mt={2}>
                  <UnorderedList color={'#1F4CED'}>
                    <ListItem fontSize="1rem">
                      <Link isExternal href="https://trisa.io/getting-started-with-trisa/">
                        Learn How TRISA Works
                      </Link>
                    </ListItem>
                    <ListItem fontSize="1rem">
                      <Link isExternal href="https://intervasp.org/">
                        What is IVMS101?
                      </Link>
                    </ListItem>
                  </UnorderedList>
                </chakra.dd>
              </GridItem>
            </Box>
          </Box>
          <Stack direction={['column', 'row']} pt={20} mx={10} justifyContent={'center'}>
            {/* <Box>
              <Button
                bg={colors.system.blue}
                color={'white'}
                _hover={{
                  bg: '#10aaed'
                }}
                _focus={{
                  borderColor: 'transparent'
                }}>
                Download PDF
              </Button>
            </Box> */}
            <Box pb={16}>
              <Button
                bg={colors.system.blue}
                px={8}
                color={'white'}
                as="a"
                href="/certificate/registration"
                _hover={{
                  bg: '#10aaed'
                }}
                _focus={{
                  borderColor: 'transparent'
                }}>
                Start Registration Process
              </Button>
            </Box>
          </Stack>
        </Stack>
      </Container>
    </Flex>
  );
}
