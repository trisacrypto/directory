import { Text, Stack } from '@chakra-ui/react';

import * as Sentry from '@sentry/react';

import AttentionAlert, { AttentionResponseType } from './AttentionAlert';
import { t } from '@lingui/macro';
import useFetchAttention from 'hooks/useFetchAttention';

export type NeedsAttentionProps = {
  text: string;
  buttonText: string;
  onClick?: (ev?: any) => void;
  loading?: boolean;
  error?: any;
  data?: Array<AttentionResponseType>;
};

const NeedsAttention = ({ buttonText, onClick }: NeedsAttentionProps) => {
  const { attentionResponse } = useFetchAttention();

  if (
    attentionResponse &&
    attentionResponse.messages &&
    !Object.keys(attentionResponse.messages).length
  ) {
    return null;
  }

  return (
    <Sentry.ErrorBoundary
      fallback={
        <Text
          color={'red'}
          textAlign={'center'}
          pt={20}>{t`An error has occurred to load attention data`}</Text>
      }>
      <Stack minHeight={67} data-cy="needs-attention">
        {attentionResponse?.messages?.map((item: AttentionResponseType, key: any) => (
          <AttentionAlert
            key={key}
            action={item?.action}
            severity={item?.severity}
            message={item?.message}
            onClick={onClick}
            buttonText={buttonText}
          />
        ))}
      </Stack>
    </Sentry.ErrorBoundary>
  );
};

export default NeedsAttention;
