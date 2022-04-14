import _ from 'lodash';
export const findStepKey = (steps: any, key: number) =>
  steps.filter((step: any) => step.key === key);

export const isValidUuid = (str: string) => {
  // Regular expression to check if string is a valid UUID
  const regexExp =
    /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
  return regexExp.test(str);
};

export const getStepData = (steps: any, key: number): TStep | undefined => {
  const s = findStepKey(steps, key);
  if (s && s?.length === 1) {
    return s[0].data;
  }
  return undefined;
};

export const getStepStatus = (steps: any, key: number): StepStatus | undefined => {
  const s = findStepKey(steps, key);
  if (s && s?.length === 1) {
    return s[0].status;
  }
  return undefined;
};

export const hasStepError = (steps: any): boolean => {
  const s = steps.filter((step: any) => step.status === 'error');
  return s.length > 0;
};

export const getStepDatas = (steps: any) => {
  const s = steps
    ?.map((step: any) => step.data)
    .reduce((acc: any, cur: any) => ({ ...acc, ...cur }), {});

  return { ...s };
};

export const getValueByPathname = (obj: Record<string, any>, path: string) => {
  return _.get(obj, path);
};

export const getDomain = (url: string | URL) => {
  try {
    const _url = new URL(url);
    return _url?.hostname?.replace('www.', '');
  } catch (error) {
    console.error('[error]', error);
    return null;
  }
};
