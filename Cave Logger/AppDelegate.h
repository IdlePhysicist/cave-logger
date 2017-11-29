//
//  AppDelegate.h
//  Cave Logger
//
//  Created by Eoghan Conlon O'Neill on 29/11/2017.
//  Copyright (c) 2017 Eoghan Conlon O'Neill. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>

@property (readonly, strong, nonatomic) NSPersistentStoreCoordinator *persistentStoreCoordinator;
@property (readonly, strong, nonatomic) NSManagedObjectModel *managedObjectModel;
@property (readonly, strong, nonatomic) NSManagedObjectContext *managedObjectContext;


@end

