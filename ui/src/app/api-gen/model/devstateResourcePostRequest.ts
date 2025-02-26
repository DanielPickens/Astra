/**
 * astra dev
 * API interface for \'astra dev\'
 *
 * The version of the OpenAPI document: 0.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


export interface DevstateResourcePostRequest { 
    /**
     * Name of the resource
     */
    name?: string;
    inlined?: string;
    uri?: string;
    deployByDefault?: DevstateResourcePostRequest.DeployByDefaultEnum;
}
export namespace DevstateResourcePostRequest {
    export type DeployByDefaultEnum = 'never' | 'undefined' | 'always';
    export const DeployByDefaultEnum = {
        Never: 'never' as DeployByDefaultEnum,
        Undefined: 'undefined' as DeployByDefaultEnum,
        Always: 'always' as DeployByDefaultEnum
    };
}


