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


export interface DevstateVolumePostRequest { 
    /**
     * Name of the volume
     */
    name?: string;
    /**
     * Minimal size of the volume
     */
    size?: string;
    /**
     * True if the Volume is Ephemeral
     */
    ephemeral?: boolean;
}

