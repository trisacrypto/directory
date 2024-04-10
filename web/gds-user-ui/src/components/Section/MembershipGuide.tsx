import { Flex, Heading, Stack, VStack, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import MembershipGuideCard from 'components/MembershipGuideCard';
import React from 'react';
import LandingLayout from 'layouts/LandingLayout';
import LandingBanner from 'components/Banner/LandingBanner';

const MembershipGuideText = [
  {
    stepNumber: 1,
    header: t`create your account`,
    description: t`Create your TRISA account with your VASP email address. Add collaborators in your organization.`,
    buttonText: t`Create Account`,
    link: '/auth/register'
  },
  {
    stepNumber: 2,
    header: t`complete VASP verification`,
    description: t`Complete the multi-part TRISA verification form and due diligence process. Once approved, gain access to the Testnet and MainNet.`,
    buttonText: t`Learn More`,
    link: '/getting-started'
  },
  {
    stepNumber: 3,
    header: t`Integrate and Comply`,
    description: t`Set up your TRISA node or integrate with a 3rd-party Travel Rule solution. Complete testing and move to production.`,
    buttonText: t`Learn More`,
    link: '/comply'
  }
];

const MembershipGuide = () => {
  return (
    <>
      <LandingLayout>
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
                <Trans id="Welcome to TRISA’s network of Certified VASPs.">
                  Welcome to TRISA’s network of Certified VASPs.
                </Trans>
              </Heading>
              <Text as="p" mt={2}>
                <Trans id="Learn about the three-step process to become a member and verified VASP.">
                  Learn about the three-step process to become a member and verified VASP.
                </Trans>
              </Text>
              <Text as="p" mt={2}>
                <Trans id="Create your account today.">Create your account today.</Trans>
              </Text>
            </VStack>
          </Stack>
        </Flex>
        <Stack
          justifyContent={'center'}
          alignContent={'center'}
          alignItems={['center', null, 'stretch']}
          spacing={10}
          direction={['column', null, 'row']}
          py={'2rem'}
          flexGrow={1}
          marginY={{ base: '2rem!important', lg: '3rem!important' }}>
          {MembershipGuideText.map(({ stepNumber, header, description, buttonText, link }) => (
            <React.Fragment key={stepNumber}>
              <MembershipGuideCard
                stepNumber={stepNumber}
                header={header}
                description={description}
                buttonText={buttonText}
                link={link}
              />
            </React.Fragment>
          ))}
        </Stack>
      </LandingLayout>
    </>
  );
};

export default MembershipGuide;
