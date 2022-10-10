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
  VStack,
  UnorderedList,
  ListItem,
  Button,
  GridItem,
  useMediaQuery
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';
import { Trans } from '@lingui/react';
import { t } from '@lingui/macro';
import { Link as RouterLink } from 'react-router-dom';

type Props = StyleProps &
  FlexProps & {
    children: React.ReactNode;
    title?: string;
    colSpan: number;
    dataContent?: string;
  };

const Line: React.FC<Props> = ({ children, colSpan, title, dataContent, ...rest }: any) => {
  const [isMobile] = useMediaQuery('(max-width: 768px)');
  return (
    <GridItem colSpan={colSpan}>
      <chakra.dt fontSize="lg" fontWeight="medium" lineHeight="6" mt={2} {...rest}>
        {title}
      </chakra.dt>
      <chakra.dd
        mt={2}
        data-content={dataContent}
        mb={{ base: 4, sm: 0 }}
        fontSize={{ base: '1rem' }}
        padding={[2, 2, 0]}
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
            <Text fontSize={'1rem'} mt={3}>
              <Trans id="TRISA members must complete a comprehensive multi-part verification form and due diligence process. Once verified, TRISA will issue TestNet and MainNet certificates for secure Travel Rule compliance.">
                TRISA members must complete a comprehensive multi-part verification form and due
                diligence process. Once verified, TRISA will issue TestNet and MainNet certificates
                for secure Travel Rule compliance.
              </Trans>
            </Text>
          </Box>
          <Box bg={'#E5EDF1'} p={6}>
            <Text fontSize={'1rem'} color={'black'}>
              <Trans id="TRISA’s verification form includes five sections and may require information from several parties in your organization.">
                TRISA’s verification form includes five sections and may require information from
                several parties in your organization.
              </Trans>
            </Text>
          </Box>
          <Box mx={20} pt={'10px'}>
            <Box
              display={{ md: 'grid' }}
              gridTemplateColumns={{ md: 'repeat(5,1fr)' }}
              color={'black'}
              gap={'10px'}
              alignItems="center">
              <>
                <Line title={t`Sections & Details`} fontWeight={'bold'} colSpan={3}>
                  {''}
                </Line>
                <Line title={t`Who to Ask`} fontWeight={'bold'} colSpan={2}>
                  {''}
                </Line>
              </>

              <Line title={t`1 Basic Details`} fontWeight={'bold'} colSpan={3}>
                <Trans id="Information about the VASP such as website, incorporation date, business and VASP category.">
                  Information about the VASP such as website, incorporation date, business and VASP
                  category.
                </Trans>
              </Line>

              <Line colSpan={2} dataContent={t`Who to ask`}>
                <Trans id="Business or Compliance Office">Business or Compliance Office</Trans>
              </Line>

              <Line title={t`2 Legal Person`} fontWeight={'bold'} colSpan={3}>
                <Trans id="Information that identifies your organization as a Legal Person. This section represents the IVMS 101 data structure for legal persons and is strongly suggested for use as KYC information exchanged in TRISA transfers.">
                  Information that identifies your organization as a Legal Person. This section
                  represents the{' '}
                  <Link color="link" href="https://intervasp.org/" isExternal>
                    IVMS 101
                  </Link>{' '}
                  data structure for legal persons and is strongly suggested for use as KYC (Know
                  your Counterparty) information exchanged in TRISA transfers.
                </Trans>
              </Line>
              <Line colSpan={2} dataContent={t`Who to ask`}>
                <Trans id="Business or Compliance Office">Business or Compliance Office</Trans>
              </Line>

              <Line title={t`3 Contacts`} fontWeight={'bold'} colSpan={3}>
                <Trans id="Contact information for representatives of your organization. Contacts include Technical, Legal/Compliance, Administrative, and Billing persons.">
                  Contact information for representatives of your organization. Contacts include
                  Technical, Legal/Compliance, Administrative, and Billing persons.
                </Trans>
              </Line>
              <Line colSpan={2} dataContent={t`Who to ask`}>
                <Trans id="Business or Compliance Office">Business or Compliance Office</Trans>
              </Line>

              <Line title={t`4 TRISA Implementation`} fontWeight={'bold'} colSpan={3}>
                <Trans id="Technical information about your endpoint for certificate issuance. Each VASP is required to establish a TRISA endpoint for inter-VASP communication.">
                  Technical information about your endpoint for certificate issuance. Each VASP is
                  required to establish a TRISA endpoint for inter-VASP communication.
                </Trans>
              </Line>
              <Line colSpan={2} dataContent={t`Who to ask`}>
                <Trans id="Technical Officer">Technical Officer</Trans>
              </Line>
              <Line title={t`5 TRIXO Questionnaire`} fontWeight={'bold'} colSpan={3}>
                <Trans id="Information to ensure that required compliance information exchanges are conducted correctly and safely. This includes information about jurisdiction and national regulator, Customer Due Diligence(CDD) and Travel Rule policies, and data protection policies.">
                  Information to ensure that required compliance information exchanges are conducted
                  correctly and safely. This includes information about jurisdiction and national
                  regulator, Customer Due Diligence(CDD) and Travel Rule policies, and data
                  protection policies.
                </Trans>
              </Line>
              <Line colSpan={2} dataContent={t`Who to ask`}>
                <Trans id="Compliance Officer">Compliance Officer</Trans>
              </Line>
            </Box>
            <Box
              display={{ md: 'grid' }}
              gridTemplateColumns={{ md: 'repeat(5,1fr)' }}
              color={'black'}
              gridColumnGap={10}
              gridRowGap={10}>
              <GridItem colSpan={3} bg={'#E5EDF1'} mt={5} p={6}>
                <chakra.dt fontSize="lg" fontWeight="bold" lineHeight="6">
                  <Trans id="Final Confirmation">Final Confirmation</Trans>
                </chakra.dt>
                <chakra.dd mt={2} fontSize="1rem">
                  <Trans id="For MainNet certificate requests, a member of TRISA’s verification team will review your submission and conduct a final due diligence phone call for physical verification. When physical verification is complete, TRISA will issue MainNet certificates. Requests for TestNet certificates do not require physical verification.">
                    For MainNet certificate requests, a member of TRISA’s verification team will
                    review your submission and conduct a final due diligence phone call for physical
                    verification. When physical verification is complete, TRISA will issue MainNet
                    certificates. Requests for TestNet certificates do not require physical
                    verification.
                  </Trans>
                </chakra.dd>
              </GridItem>
              <GridItem colSpan={2} bg={'#E5EDF1'} mt={5} p={6}>
                <chakra.dt fontSize="lg" fontWeight="bold" lineHeight="6">
                  <Trans id="Need to Learn More?">Need to Learn More?</Trans>
                </chakra.dt>
                <chakra.dd mt={2}>
                  <UnorderedList color={'#1F4CED'}>
                    <ListItem fontSize="1rem">
                      <Link isExternal href="https://trisa.io/getting-started-with-trisa/">
                        <Trans id="Learn How TRISA Works">Learn How TRISA Works</Trans>
                      </Link>
                    </ListItem>
                    <ListItem fontSize="1rem">
                      <Link isExternal href="https://intervasp.org/">
                        <Trans id="What is IVMS101?">What is IVMS101?</Trans>
                      </Link>
                    </ListItem>
                  </UnorderedList>
                </chakra.dd>
              </GridItem>
            </Box>
          </Box>
          <Stack
            direction={['column', 'row']}
            pb={20}
            pt={5}
            justifyContent={'center'}
            textAlign={'center'}>
            <VStack fontSize={'md'}>
              <RouterLink to={'/auth/register'}>
                <Button
                  bg={colors.system.blue}
                  color={'white'}
                  minWidth={'300px'}
                  _hover={{
                    bg: '#10aaed'
                  }}
                  _focus={{
                    borderColor: 'transparent'
                  }}>
                  <Trans id="Create account">Create account</Trans>
                </Button>
              </RouterLink>
              <Text textAlign="center">
                <Trans id="Already have an account?">Already have an account?</Trans>{' '}
                <RouterLink to={'/auth/login'}>
                  <Link color={colors.system.cyan}>
                    {' '}
                    <Trans id="Log in.">Log in.</Trans>
                  </Link>
                </RouterLink>
              </Text>
            </VStack>
          </Stack>
        </Stack>
      </Container>
    </Flex>
  );
}
