import { Meta, Story } from "@storybook/react";
import CertificateSection from ".";

type CertificateSectionProps = {
  step: number;
  title?: string;
  description?: string;
  isSaved?: boolean;
  isSubmitted?: boolean;
};

export default {
  title: "components/CertificateSection",
  component: CertificateSection,
} as Meta;

const Template: Story<CertificateSectionProps> = (args) => (
  <CertificateSection {...args} />
);

export const Default = Template.bind({});
Default.args = {
  step: 2,
};
