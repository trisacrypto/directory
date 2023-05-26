import React from 'react';
import { SimpleDashboardLayout } from 'layouts';
import { Box, Heading, VStack, Text, Link, Stack, useColorModeValue, Flex } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/RegistrationForm/CertificateRegistrationForm';

import { userSelector } from 'modules/auth/login/user.slice';

import { useSelector } from 'react-redux';

import HomeButton from 'components/ui/HomeButton';

import { Trans } from '@lingui/react';

const Certificate: React.FC = () => {
  const textColor = useColorModeValue('black', '#EDF2F7');
  const backgroundColor = useColorModeValue('white', '#171923');

  const { isLoggedIn } = useSelector(userSelector);

  return (
    <SimpleDashboardLayout>
      <>
        <Flex justifyContent={'space-between'}>
          <Heading size="lg" mb="24px" className="heading">
            <Trans id="Certificate Registration">Certificate Registration</Trans>
          </Heading>
          <Box>{!isLoggedIn && <HomeButton link={'/'} />}</Box>
        </Flex>
        <Stack my={3}>
          <Card maxW="100%" bg={backgroundColor} color={textColor}>
            <Card.Body>
              <Text>
                <Trans id="This multi-section form is an important step in the registration and certificate issuance process. The information you provide will be used to verify the legal entity that you represent and, where appropriate, will be available to verified TRISA members to facilitate compliance decisions. If you need guidance, see the">
                  This multi-section form is an important step in the registration and certificate
                  issuance process. The information you provide will be used to verify the legal
                  entity that you represent and, where appropriate, will be available to verified
                  TRISA members to facilitate compliance decisions. If you need guidance, see the
                </Trans>{' '}
                <Link isExternal href="/getting-started" color={'link'} fontWeight={'bold'}>
                  <Trans id="Getting Started Help Guide">Getting Started Help Guide</Trans>.{' '}
                </Link>
              </Text>
              <Text pt={4}>
                <Trans id="To assist in completing the registration form, the form is divided into multiple sections">
                  To assist in completing the registration form, the form is divided into multiple
                  sections
                </Trans>
                .{' '}
                <Text as={'span'} fontWeight={'bold'}>
                  <Trans id="No information is sent until you complete Section 6 - Review & Submit">
                    No information is sent until you complete Section 6 - Review & Submit
                  </Trans>
                  .{' '}
                </Text>
              </Text>
            </Card.Body>
          </Card>
        </Stack>

        <>
          <VStack spacing={3}>
            <Box width={'100%'}>
              <TestNetCertificateProgressBar />
            </Box>
            <Stack
              width="100%"
              direction={'row'}
              spacing={8}
              justifyContent={'center'}
              py={6}
              wrap="wrap"
              rowGap={2}></Stack>
          </VStack>
        </>
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
