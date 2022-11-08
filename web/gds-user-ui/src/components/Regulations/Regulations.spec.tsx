import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, waitFor } from 'utils/test-utils';
import Regulations from '.';

describe('<Regulations />', () => {
  beforeAll(async () => {
    await waitFor(() => {
      dynamicActivate('en');
    });
  });

  it('should should match snapshot', () => {
    const { container } = render(<Regulations name="" />);

    expect(container).toMatchInlineSnapshot(`
      .emotion-0 {
        display: -webkit-box;
        display: -webkit-flex;
        display: -ms-flexbox;
        display: flex;
        -webkit-align-items: start;
        -webkit-box-align: start;
        -ms-flex-align: start;
        align-items: start;
        -webkit-flex-direction: column;
        -ms-flex-direction: column;
        flex-direction: column;
        width: 100%;
      }

      .emotion-0>*:not(style)~*:not(style) {
        margin-top: 0.5rem;
        -webkit-margin-end: 0px;
        margin-inline-end: 0px;
        margin-bottom: 0px;
        -webkit-margin-start: 0px;
        margin-inline-start: 0px;
      }

      .emotion-1 {
        display: -webkit-inline-box;
        display: -webkit-inline-flex;
        display: -ms-inline-flexbox;
        display: inline-flex;
        -webkit-appearance: none;
        -moz-appearance: none;
        -ms-appearance: none;
        appearance: none;
        -webkit-align-items: center;
        -webkit-box-align: center;
        -ms-flex-align: center;
        align-items: center;
        -webkit-box-pack: center;
        -ms-flex-pack: center;
        -webkit-justify-content: center;
        justify-content: center;
        -webkit-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
        position: relative;
        white-space: nowrap;
        vertical-align: middle;
        outline: 2px solid transparent;
        outline-offset: 2px;
        width: auto;
        line-height: 1.2;
        border-radius: 5px;
        font-weight: var(--chakra-fontWeights-semibold);
        transition-property: var(--chakra-transition-property-common);
        transition-duration: var(--chakra-transition-duration-normal);
        height: var(--chakra-sizes-10);
        min-width: var(--chakra-sizes-10);
        font-size: var(--chakra-fontSizes-md);
        -webkit-padding-start: var(--chakra-space-4);
        padding-inline-start: var(--chakra-space-4);
        -webkit-padding-end: var(--chakra-space-4);
        padding-inline-end: var(--chakra-space-4);
        background: var(--chakra-colors-gray-100);
      }

      .emotion-1:focus,
      .emotion-1[data-focus] {
        box-shadow: var(--chakra-shadows-outline);
      }

      .emotion-1[disabled],
      .emotion-1[aria-disabled=true],
      .emotion-1[data-disabled] {
        opacity: 0.4;
        cursor: not-allowed;
        box-shadow: var(--chakra-shadows-none);
      }

      .emotion-1:hover,
      .emotion-1[data-hover] {
        background: var(--chakra-colors-gray-200);
      }

      .emotion-1:hover[disabled],
      .emotion-1[data-hover][disabled],
      .emotion-1:hover[aria-disabled=true],
      .emotion-1[data-hover][aria-disabled=true],
      .emotion-1:hover[data-disabled],
      .emotion-1[data-hover][data-disabled] {
        background: var(--chakra-colors-gray-100);
      }

      .emotion-1:active,
      .emotion-1[data-active] {
        background: var(--chakra-colors-gray-300);
      }

      <div>
        <div
          class="chakra-stack emotion-0"
        >
          <button
            class="chakra-button emotion-1"
            type="button"
          >
            Add Regulation
          </button>
        </div>
      </div>
    `);
  });
});
