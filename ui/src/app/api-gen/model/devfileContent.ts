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
import { Container } from './container';
import { Command } from './command';
import { Events } from './events';
import { Volume } from './volume';
import { Metadata } from './metadata';
import { Resource } from './resource';
import { Image } from './image';


export interface DevfileContent { 
    content: string;
    version: string;
    commands: Array<Command>;
    containers: Array<Container>;
    images: Array<Image>;
    resources: Array<Resource>;
    volumes: Array<Volume>;
    events: Events;
    metadata: Metadata;
}

