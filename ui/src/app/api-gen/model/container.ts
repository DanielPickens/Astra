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
import { Endpoint } from './endpoint';
import { VolumeMount } from './volumeMount';
import { Env } from './env';
import { Annotation } from './annotation';


export interface Container { 
    name: string;
    image: string;
    command: Array<string>;
    args: Array<string>;
    memoryRequest: string;
    memoryLimit: string;
    cpuRequest: string;
    cpuLimit: string;
    volumeMounts: Array<VolumeMount>;
    annotation: Annotation;
    endpoints: Array<Endpoint>;
    env: Array<Env>;
    configureSources: boolean;
    mountSources: boolean;
    sourceMapping: string;
}

