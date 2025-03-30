import React from 'react';
import { Stack, Text, HStack, Link } from '@chakra-ui/react';
import LandingLayout from 'layouts/LandingLayout';
import { Link as RouteLink } from 'react-router-dom';
import SuccessSvg from 'assets/successSvg.svg';
import { Trans } from '@lingui/react';
import { CircleChevronRight } from 'akar-icons';
import CkLazyLoadImage from 'components/LazyImage';

const VerifyPage: React.FC = () => {
  return (
    <LandingLayout>
      <Stack direction="row" spacing={8} mx={'auto'} maxW={'3xl'} width={'100%'}>
        <Stack pt={5}>
          <CkLazyLoadImage src={SuccessSvg} width="100px" mx="auto" py={10} />
          <Text fontSize="xl" fontWeight="bold">
            <Trans id="Thank you for creating your TRISA account.">
              Thank you for creating your TRISA account.{' '}
            </Trans>
          </Text>

          <Text py={5}>
            <Trans id="Thank you for creating your TRISA account. You must verify your email address. An email with verification instructions has been sent to your email address. Please complete the email verification process in 24 hours. The email verification link will expire in 24 hours.">
              Thank you for creating your TRISA account. You must verify your email address. An
              email with verification instructions has been sent to your email address. Please
              complete the email verification process in 24 hours. The email verification link will
              expire in 24 hours.
            </Trans>
          </Text>

          <HStack spacing={4} py={3} color={'#55ACD8'}>
            <CircleChevronRight strokeWidth={2} size={36} />
            <Text as={'span'} ml={3}>
              <RouteLink to={'/auth/login'}>
                <Link _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                  <Trans id="Log in to TRISA’s Global Directory Service">
                    Log in to TRISA’s Global Directory Service{' '}
                  </Trans>
                </Link>
              </RouteLink>
            </Text>
          </HStack>
          <HStack spacing={4} pb={3} color={'#55ACD8'}>
            <CircleChevronRight strokeWidth={2} size={36} />
            <Text as={'span'} ml={3}>
              <RouteLink to={'/getting-started'}>
                <Link _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                  <Trans id="Visit Getting Started with TRISA">
                    Visit Getting Started with TRISA{' '}
                  </Trans>
                </Link>
              </RouteLink>
            </Text>
          </HStack>
          <HStack spacing={4} pb={5} color={'#55ACD8'}>
            <CircleChevronRight strokeWidth={2} size={36} />
            <Text as={'span'} ml={3}>
              <RouteLink to={'/'}>
                <Link _active={{ outline: 'none' }} _focus={{ outline: 'none' }}>
                  <Trans id="Return to trisa.directory">Return to trisa.directory </Trans>
                </Link>
              </RouteLink>
            </Text>
          </HStack>
        </Stack>
      </Stack>
    </LandingLayout>
  );
};

export default VerifyPage;
